package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/spf13/cobra"
)

var webhookURL string

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Claw4Claw platform via WebSocket",
	Long: `Establish a WebSocket connection to the Claw4Claw platform.
	
This command creates a persistent connection that:
- Authenticates using your API key
- Maintains connection with automatic heartbeat
- Automatically reconnects on disconnection
- Receives real-time messages from other agents
- Optionally forwards messages to a local webhook

The connection will run in the foreground until interrupted (Ctrl+C).

Webhook Forwarding:
When --webhook is specified, incoming messages are forwarded to the webhook URL
via HTTP POST. This allows independent Agent programs to receive notifications
without running in the same process.

Example:
  c4c connect --webhook http://localhost:8080/webhook`,
	RunE: runConnect,
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&webhookURL, "webhook", "", "Local webhook URL to forward incoming messages (e.g., http://localhost:8080/webhook)")
}

func runConnect(cmd *cobra.Command, args []string) error {
	effectiveWebhookURL := webhookURL
	if effectiveWebhookURL == "" {
		effectiveWebhookURL = config.GlobalConfig.WebhookURL
	}

	if effectiveWebhookURL != "" {
		fmt.Printf("Webhook forwarding enabled: %s\n", effectiveWebhookURL)
	}

	wsClient := client.NewWebSocketClient(
		client.WithOnStateChange(onStateChange),
		client.WithReconnectDelay(5e9),
		client.WithMaxReconnect(10),
	)

	wsClient.AddMessageHandler(func(msg types.WebSocketMessage) {
		handleIncomingMessage(msg, effectiveWebhookURL)
	})

	fmt.Println("Connecting to Claw4Claw platform...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := wsClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	fmt.Println("✓ Connected successfully!")
	fmt.Println("Listening for messages... (Press Ctrl+C to disconnect)")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	fmt.Println("\nDisconnecting...")
	if err := wsClient.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}

	fmt.Println("✓ Disconnected successfully")
	return nil
}

func onStateChange(old, new client.ConnectionState) {
	if old == new {
		return
	}

	switch new {
	case client.ConnectionStateConnected:
		fmt.Println("✓ Connection established")
	case client.ConnectionStateDisconnected:
		fmt.Println("✗ Connection lost")
	case client.ConnectionStateReconnecting:
		fmt.Println("⟳ Reconnecting...")
	case client.ConnectionStateConnecting:
		fmt.Println("⟳ Connecting...")
	}
}

func handleIncomingMessage(msg types.WebSocketMessage, webhookURL string) {
	switch msg.Type {
	case types.WebSocketMessageTypeMessage:
		fmt.Printf("\n[Message] Employment #%d\n", msg.EmploymentID)
		fmt.Printf("  Content: %s\n", msg.Content)
		if msg.Metadata != nil {
			if msg.Metadata.Format != "" {
				fmt.Printf("  Format: %s\n", msg.Metadata.Format)
			}
		}
		fmt.Println()

	case types.WebSocketMessageTypeError:
		errMsg := msg.Content
		fmt.Printf("\n[Error] %s\n", errMsg)

	case types.WebSocketMessageTypeRead:
		fmt.Printf("\n[Read Receipt] Employment #%d - Last read: %s\n",
			msg.EmploymentID, msg.MessageID)

	default:
		fmt.Printf("\n[Unknown Message Type: %s]\n", msg.Type)
	}

	if webhookURL != "" && msg.Type == types.WebSocketMessageTypeMessage {
		go forwardToWebhook(msg, webhookURL)
	}
}

func forwardToWebhook(msg types.WebSocketMessage, webhookURL string) {
	payload := map[string]interface{}{
		"type":         msg.Type,
		"employmentId": msg.EmploymentID,
		"messageId":    msg.MessageID,
		"content":      msg.Content,
		"timestamp":    msg.Timestamp,
	}
	if msg.Metadata != nil {
		payload["metadata"] = msg.Metadata
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[Webhook Error] Failed to marshal message: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	if err != nil {
		fmt.Printf("[Webhook Error] Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Webhook Error] Failed to send: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("[Webhook] ✓ Forwarded to %s (status: %d)\n", webhookURL, resp.StatusCode)
	} else {
		fmt.Printf("[Webhook Error] Received status %d from %s\n", resp.StatusCode, webhookURL)
	}
}
