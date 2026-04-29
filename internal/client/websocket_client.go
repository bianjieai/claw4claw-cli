package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionState string

const (
	ConnectionStateDisconnected ConnectionState = "disconnected"
	ConnectionStateConnecting   ConnectionState = "connecting"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateReconnecting ConnectionState = "reconnecting"
)

type MessageHandler func(msg types.WebSocketMessage)
type NotificationHandler func(msg types.WebSocketNotificationMessage)

type WebSocketClient struct {
	conn                 *websocket.Conn
	connMutex            sync.RWMutex
	writeMutex           sync.Mutex // 写操作互斥锁
	state                ConnectionState
	stateMutex           sync.RWMutex
	apiKey               string
	wsEndpoint           string
	stopChan             chan struct{}
	stopOnce             sync.Once
	messageHandlers      []MessageHandler
	notificationHandlers []NotificationHandler
	handlerMutex         sync.RWMutex
	reconnectDelay       time.Duration
	maxReconnect         int
	reconnectCount       int
	sendChan             chan types.WebSocketMessage
	pingTicker           *time.Ticker
	pongTimeout          time.Duration
	lastPongTime         time.Time
	lastPongMutex        sync.RWMutex
	onStateChange        func(old, new ConnectionState)
}

type WebSocketClientOption func(*WebSocketClient)

func WithReconnectDelay(delay time.Duration) WebSocketClientOption {
	return func(c *WebSocketClient) {
		c.reconnectDelay = delay
	}
}

func WithMaxReconnect(max int) WebSocketClientOption {
	return func(c *WebSocketClient) {
		c.maxReconnect = max
	}
}

func WithPongTimeout(timeout time.Duration) WebSocketClientOption {
	return func(c *WebSocketClient) {
		c.pongTimeout = timeout
	}
}

func WithOnStateChange(handler func(old, new ConnectionState)) WebSocketClientOption {
	return func(c *WebSocketClient) {
		c.onStateChange = handler
	}
}

func NewWebSocketClient(opts ...WebSocketClientOption) *WebSocketClient {
	wsEndpoint := convertHTTPToWS(config.GlobalConfig.APIEndpoint)

	client := &WebSocketClient{
		apiKey:               config.GlobalConfig.APIToken,
		wsEndpoint:           wsEndpoint + "/ws",
		state:                ConnectionStateDisconnected,
		stopChan:             make(chan struct{}),
		messageHandlers:      make([]MessageHandler, 0),
		notificationHandlers: make([]NotificationHandler, 0),
		reconnectDelay:       5 * time.Second,
		maxReconnect:         10,
		sendChan:             make(chan types.WebSocketMessage, 100),
		pongTimeout:          120 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func convertHTTPToWS(endpoint string) string {
	if strings.HasPrefix(endpoint, "https://") {
		return strings.Replace(endpoint, "https://", "wss://", 1)
	}
	if strings.HasPrefix(endpoint, "http://") {
		return strings.Replace(endpoint, "http://", "ws://", 1)
	}
	return endpoint
}

func (c *WebSocketClient) Connect(ctx context.Context) error {
	c.setState(ConnectionStateConnecting)

	u, err := url.Parse(c.wsEndpoint)
	if err != nil {
		c.setState(ConnectionStateDisconnected)
		return fmt.Errorf("failed to parse websocket endpoint: %w", err)
	}

	q := u.Query()
	q.Set("api_key", c.apiKey)
	u.RawQuery = q.Encode()

	headers := http.Header{}
	headers.Set("X-API-Key", c.apiKey)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		HandshakeTimeout: 10 * time.Second,
	}

	conn, resp, err := dialer.Dial(u.String(), headers)
	if err != nil {
		c.setState(ConnectionStateDisconnected)
		if resp != nil && resp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: invalid API key")
		}
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	c.setConnection(conn)
	c.setState(ConnectionStateConnected)
	c.reconnectCount = 0

	conn.SetPongHandler(func(appData string) error {
		c.lastPongMutex.Lock()
		c.lastPongTime = time.Now()
		c.lastPongMutex.Unlock()
		return nil
	})

	go c.readMessages()
	go c.writeMessages()
	go c.heartbeat()

	return nil
}

func (c *WebSocketClient) Disconnect() error {
	c.stopOnce.Do(func() {
		close(c.stopChan)
		if c.pingTicker != nil {
			c.pingTicker.Stop()
		}
		c.closeConnection()
		c.setState(ConnectionStateDisconnected)
	})
	return nil
}

func (c *WebSocketClient) IsConnected() bool {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	return c.state == ConnectionStateConnected
}

func (c *WebSocketClient) GetState() ConnectionState {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()
	return c.state
}

func (c *WebSocketClient) SendMessage(msg types.WebSocketMessage) error {
	if !c.IsConnected() {
		return fmt.Errorf("websocket is not connected")
	}

	if msg.MessageID == "" {
		msg.MessageID = generateMessageID()
	}
	if msg.Timestamp == "" {
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	select {
	case c.sendChan <- msg:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout, message not delivered")
	}
}

func (c *WebSocketClient) AddMessageHandler(handler MessageHandler) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()
	c.messageHandlers = append(c.messageHandlers, handler)
}

func (c *WebSocketClient) AddNotificationHandler(handler NotificationHandler) {
	c.handlerMutex.Lock()
	defer c.handlerMutex.Unlock()
	c.notificationHandlers = append(c.notificationHandlers, handler)
}

func (c *WebSocketClient) readMessages() {
	for {
		select {
		case <-c.stopChan:
			return
		default:
		}

		conn := c.getConnection()
		if conn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.handleDisconnect()
			}
			return
		}

		var raw json.RawMessage
		if err := json.Unmarshal(message, &raw); err != nil {
			continue
		}

		var typeMsg struct {
			Type types.WebSocketMessageType `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeMsg); err != nil {
			continue
		}

		switch typeMsg.Type {
		case types.WebSocketMessageTypePing:
			var pingMsg types.WebSocketPingMessage
			if err := json.Unmarshal(raw, &pingMsg); err != nil {
				continue
			}
			c.handlePing(pingMsg)
		case types.WebSocketMessageTypeNotification:
			var notif types.WebSocketNotificationMessage
			if err := json.Unmarshal(raw, &notif); err != nil {
				continue
			}
			c.handleNotification(notif)
		default:
			var msg types.WebSocketMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}
			c.handleMessage(msg)
		}
	}
}

func (c *WebSocketClient) writeMessages() {
	for {
		select {
		case <-c.stopChan:
			return
		case msg := <-c.sendChan:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			if err := c.writeMessage(websocket.TextMessage, data); err != nil {
				c.handleDisconnect()
			}
		}
	}
}

func (c *WebSocketClient) heartbeat() {
	c.pingTicker = time.NewTicker(60 * time.Second)
	defer c.pingTicker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-c.pingTicker.C:
			if err := c.writeMessage(websocket.PingMessage, nil); err != nil {
				c.handleDisconnect()
				continue
			}

			c.lastPongMutex.RLock()
			lastPong := c.lastPongTime
			c.lastPongMutex.RUnlock()

			if !lastPong.IsZero() && time.Since(lastPong) > c.pongTimeout {
				c.handleDisconnect()
			}
		}
	}
}

func (c *WebSocketClient) handleMessage(msg types.WebSocketMessage) {
	c.handlerMutex.RLock()
	handlers := make([]MessageHandler, len(c.messageHandlers))
	copy(handlers, c.messageHandlers)
	c.handlerMutex.RUnlock()

	for _, handler := range handlers {
		go handler(msg)
	}
}

func (c *WebSocketClient) handleNotification(notif types.WebSocketNotificationMessage) {
	c.handlerMutex.RLock()
	handlers := make([]NotificationHandler, len(c.notificationHandlers))
	copy(handlers, c.notificationHandlers)
	c.handlerMutex.RUnlock()

	for _, handler := range handlers {
		go handler(notif)
	}
}

func (c *WebSocketClient) handlePing(ping types.WebSocketPingMessage) {
	pong := types.WebSocketPongMessage{
		Type:      types.WebSocketMessageTypePong,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	pongBytes, err := json.Marshal(pong)
	if err != nil {
		return
	}

	if err := c.writeMessage(websocket.TextMessage, pongBytes); err != nil {
		return
	}
}

func (c *WebSocketClient) handleDisconnect() {
	c.closeConnection()

	if c.reconnectCount < c.maxReconnect {
		go c.reconnect()
	} else {
		c.setState(ConnectionStateDisconnected)
	}
}

func (c *WebSocketClient) reconnect() {
	c.setState(ConnectionStateReconnecting)
	c.reconnectCount++

	time.Sleep(c.reconnectDelay)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Connect(ctx); err != nil {
		c.handleDisconnect()
	}
}

func (c *WebSocketClient) setState(state ConnectionState) {
	c.stateMutex.Lock()
	oldState := c.state
	c.state = state
	c.stateMutex.Unlock()

	if c.onStateChange != nil && oldState != state {
		c.onStateChange(oldState, state)
	}
}

func (c *WebSocketClient) setConnection(conn *websocket.Conn) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	c.conn = conn
}

func (c *WebSocketClient) getConnection() *websocket.Conn {
	c.connMutex.RLock()
	defer c.connMutex.RUnlock()
	return c.conn
}

func (c *WebSocketClient) closeConnection() {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	if c.conn != nil {
		// 使用写锁保护关闭消息
		c.writeMutex.Lock()
		_ = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.writeMutex.Unlock()
		c.conn.Close()
		c.conn = nil
	}
}

// writeMessage 线程安全的写消息方法
func (c *WebSocketClient) writeMessage(messageType int, data []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	conn := c.getConnection()
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	return conn.WriteMessage(messageType, data)
}

func generateMessageID() string {
	return "msg_" + uuid.New().String()
}
