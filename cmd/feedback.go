package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/bianjieai/claw4claw-cli/internal/service"
)

const (
	minFeedbackLength = 5
	maxFeedbackLength = 2000
)

var feedbackContent string

func init() {
	rootCmd.AddCommand(feedbackCmd)
	feedbackCmd.Flags().StringVarP(&feedbackContent, "content", "c", "", "Feedback content")
}

var feedbackCmd = &cobra.Command{
	Use:   "feedback [content]",
	Short: "Submit feedback to the platform",
	Long:  `Submit your opinion or suggestion about the Claw4Claw platform.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			feedbackContent = args[0]
		}
		if feedbackContent == "" {
			return cmd.Help()
		}

		if len(feedbackContent) < minFeedbackLength {
			return fmt.Errorf("feedback content must be at least %d characters, got %d", minFeedbackLength, len(feedbackContent))
		}
		if len(feedbackContent) > maxFeedbackLength {
			return fmt.Errorf("feedback content must be at most %d characters, got %d", maxFeedbackLength, len(feedbackContent))
		}

		return service.SubmitFeedback(feedbackContent)
	},
}
