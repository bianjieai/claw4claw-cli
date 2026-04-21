package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/spf13/cobra"
)

var (
	chatMessage     string
	chatInteractive bool
	chatHistory     bool
	chatLimit       int
)

var chatCmd = &cobra.Command{
	Use:   "chat <employment-id>",
	Short: "Chat with an employed agent",
	Long: `Send messages or interact with an agent in an employment relationship.

This command supports three modes:

1. Send a single message:
   c4c chat 123 --message "Please help me with this task"

2. Interactive chat mode:
   c4c chat 123 --interactive

3. View message history:
   c4c chat 123 --history --limit 50

The interactive mode provides a real-time chat interface where you can:
- Type messages and press Enter to send
- Type 'exit' or 'quit' to end the session
- Press Ctrl+C to disconnect`,
	Args: cobra.ExactArgs(1),
	RunE: runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&chatMessage, "message", "m", "", "Send a single message")
	chatCmd.Flags().BoolVarP(&chatInteractive, "interactive", "i", false, "Enter interactive chat mode")
	chatCmd.Flags().BoolVarP(&chatHistory, "history", "", false, "View message history")
	chatCmd.Flags().IntVarP(&chatLimit, "limit", "l", 20, "Number of messages to retrieve (for history mode)")
}

func runChat(cmd *cobra.Command, args []string) error {
	employmentID, err := strconv.ParseUint(args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid employment ID: %w", err)
	}

	switch {
	case chatHistory:
		return viewHistory(uint(employmentID))
	case chatMessage != "":
		return sendMessage(uint(employmentID), chatMessage)
	case chatInteractive:
		return startInteractiveMode(uint(employmentID))
	default:
		return fmt.Errorf("please specify one of: --message, --interactive, or --history")
	}
}

func viewHistory(employmentID uint) error {
	apiClient := client.NewAPIClient()

	params := types.GetMessagesQueryParams{
		Limit: chatLimit,
	}

	fmt.Printf("Fetching message history for employment #%d...\n\n", employmentID)

	resp, err := apiClient.GetEmploymentMessages(employmentID, params)
	if err != nil {
		return fmt.Errorf("failed to get message history: %w", err)
	}

	if len(resp.Messages) == 0 {
		fmt.Println("No messages found.")
		return nil
	}

	fmt.Printf("Messages (%d):\n", len(resp.Messages))
	fmt.Println(strings.Repeat("-", 80))

	for _, msg := range resp.Messages {
		sender := "You"
		if msg.SenderAgentID != 0 {
			sender = fmt.Sprintf("Agent #%d", msg.SenderAgentID)
		}

		timestamp := msg.CreatedAt
		if t, err := time.Parse(time.RFC3339, msg.CreatedAt); err == nil {
			timestamp = t.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("[%s] %s:\n", timestamp, sender)
		fmt.Printf("  %s\n", msg.Content)
		if msg.ReadAt != nil {
			fmt.Printf("  ✓ Read\n")
		}
		fmt.Println()
	}

	if resp.HasMore {
		fmt.Printf("... more messages available (cursor: %s)\n", resp.NextCursor)
	}

	return nil
}

func sendMessage(employmentID uint, content string) error {
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	wsClient := client.NewWebSocketClient(
		client.WithOnStateChange(onStateChange),
	)

	wsClient.AddMessageHandler(func(msg types.WebSocketMessage) {
		if msg.Type == types.WebSocketMessageTypeMessage && msg.EmploymentID == employmentID {
			fmt.Printf("\n[Reply] %s\n", msg.Content)
		} else if msg.Type == types.WebSocketMessageTypeError {
			fmt.Printf("\n[Error] %s\n", msg.Content)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := wsClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer wsClient.Disconnect()

	msg := types.WebSocketMessage{
		Type:         types.WebSocketMessageTypeMessage,
		EmploymentID: employmentID,
		Content:      content,
	}

	fmt.Printf("Sending message to employment #%d...\n", employmentID)

	if err := wsClient.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Println("✓ Message sent")

	fmt.Println("Waiting for reply... (Press Ctrl+C to exit)")
	time.Sleep(5 * time.Second)

	return nil
}

func startInteractiveMode(employmentID uint) error {
	wsClient := client.NewWebSocketClient(
		client.WithOnStateChange(onStateChange),
	)

	wsClient.AddMessageHandler(func(msg types.WebSocketMessage) {
		if msg.Type == types.WebSocketMessageTypeMessage && msg.EmploymentID == employmentID {
			timestamp := time.Now().Format("15:04:05")
			fmt.Printf("\n[%s] Agent: %s\n", timestamp, msg.Content)
			fmt.Print("> ")
		} else if msg.Type == types.WebSocketMessageTypeError {
			fmt.Printf("\n[Error] %s\n", msg.Content)
			fmt.Print("> ")
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("Connecting to employment #%d...\n", employmentID)

	if err := wsClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer wsClient.Disconnect()

	fmt.Println("✓ Connected")
	fmt.Println("Interactive chat mode. Type your message and press Enter to send.")
	fmt.Println("Type 'exit' or 'quit' to end the session.")
	fmt.Println(strings.Repeat("-", 80))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)
	inputChan := make(chan string, 1)

	go func() {
		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				close(inputChan)
				return
			}
			inputChan <- strings.TrimSpace(input)
		}
	}()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\nDisconnecting...")
			wsClient.Disconnect()
			return nil

		case input, ok := <-inputChan:
			if !ok {
				wsClient.Disconnect()
				return nil
			}

			if input == "" {
				continue
			}

			if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
				fmt.Println("Ending chat session...")
				wsClient.Disconnect()
				return nil
			}

			msg := types.WebSocketMessage{
				Type:         types.WebSocketMessageTypeMessage,
				EmploymentID: employmentID,
				Content:      input,
			}

			if err := wsClient.SendMessage(msg); err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			}
		}
	}
}
