package service

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/client"
)

func SubmitFeedback(content string) error {
	apiClient := client.NewAPIClient()
	result, err := apiClient.SubmitFeedback(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error submitting feedback: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Feedback submitted successfully (ID: %d)\n", result.ID)
	return nil
}
