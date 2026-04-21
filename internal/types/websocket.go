package types

import "time"

type WebSocketMessageType string

const (
	WebSocketMessageTypeMessage WebSocketMessageType = "message"
	WebSocketMessageTypePing    WebSocketMessageType = "ping"
	WebSocketMessageTypePong    WebSocketMessageType = "pong"
	WebSocketMessageTypeError   WebSocketMessageType = "error"
	WebSocketMessageTypeRead    WebSocketMessageType = "read"
)

type WebSocketMessage struct {
	Type        WebSocketMessageType `json:"type"`
	EmploymentID uint                `json:"employmentId,omitempty"`
	Timestamp   string               `json:"timestamp,omitempty"`
	MessageID   string               `json:"messageId,omitempty"`
	Content     string               `json:"content,omitempty"`
	Metadata    *MessageMetadata     `json:"metadata,omitempty"`
}

type MessageMetadata struct {
	Language    string       `json:"language,omitempty"`
	Format      string       `json:"format,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Type string `json:"type"`
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
	Size int64  `json:"size,omitempty"`
}

type WebSocketPingMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Timestamp string               `json:"timestamp"`
}

type WebSocketPongMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Timestamp string               `json:"timestamp"`
}

type WebSocketErrorMessage struct {
	Type    WebSocketMessageType `json:"type"`
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Details interface{}          `json:"details,omitempty"`
}

type WebSocketReadMessage struct {
	Type              WebSocketMessageType `json:"type"`
	EmploymentID      uint                `json:"employmentId"`
	Timestamp         string               `json:"timestamp"`
	MessageID         string               `json:"messageId"`
	LastReadMessageID string               `json:"lastReadMessageId"`
}

type EmploymentMessage struct {
	ID              uint                `json:"id"`
	EmploymentID    uint                `json:"employmentId"`
	SenderAgentID   uint                `json:"senderAgentId"`
	ReceiverAgentID uint                `json:"receiverAgentId"`
	Type            WebSocketMessageType `json:"type"`
	Content         string               `json:"content"`
	MessageID       string               `json:"messageId"`
	ReadAt          *time.Time          `json:"readAt,omitempty"`
	CreatedAt       string               `json:"createdAt"`
}

type EmploymentMessagesResponse struct {
	Messages   []EmploymentMessage `json:"messages"`
	HasMore    bool                `json:"hasMore"`
	NextCursor string              `json:"nextCursor"`
}

type MarkMessagesReadRequest struct {
	LastReadMessageID string `json:"lastReadMessageId"`
}

type MarkMessagesReadResponse struct {
	UnreadCount int `json:"unreadCount"`
}

type GetMessagesQueryParams struct {
	Before string
	Limit  int
}
