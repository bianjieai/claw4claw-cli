package manage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml.v3"
)

var (
	taskRole   string
	taskSearch string
	taskStatus string
	taskSortBy string

	acceptedStatus string
	acceptedOutput string

	publishTitle       string
	publishDescription string
	publishBounty      float64
	publishCategory    string
	publishDeadline    string
	publishFile        string
	publishOutput      string

	applyAgent         string
	applyMessage       string
	applyEstimatedTime string
	applyOutput        string

	submitContent     string
	submitFile        string
	submitAttachments []string
	submitNotes       string
	submitOutput      string

	acceptRating int
	acceptReview string
	acceptOutput string

	applicationsStatus string
	applicationsOutput string

	acceptApplicantMessage string
	acceptApplicantOutput  string

	cancelOutput string

	reviewOutput string
)

func init() {
	ManageCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskAcceptedCmd)
	taskCmd.AddCommand(taskPublishCmd)
	taskCmd.AddCommand(taskApplyCmd)
	taskCmd.AddCommand(taskSubmitCmd)
	taskCmd.AddCommand(taskAcceptCmd)
	taskCmd.AddCommand(taskApplicationsCmd)
	taskCmd.AddCommand(taskAcceptApplicantCmd)
	taskCmd.AddCommand(taskCancelCmd)
	taskCmd.AddCommand(taskReviewCmd)

	taskListCmd.Flags().StringVarP(&taskRole, "role", "r", "all", "Role filter (publisher/worker/all)")
	taskListCmd.Flags().StringVarP(&taskSearch, "search", "s", "", "Search keyword")
	taskListCmd.Flags().StringVar(&taskStatus, "status", "all", "Status filter (all/open/in_progress/pending_review/completed/cancelled)")
	taskListCmd.Flags().StringVar(&taskSortBy, "sort-by", "time", "Sort by (time/bounty)")

	taskAcceptedCmd.Flags().StringVar(&acceptedStatus, "status", "all", "Status filter (all/in_progress/pending_review/completed)")
	taskAcceptedCmd.Flags().StringVarP(&acceptedOutput, "output", "o", "text", "Output format (text/json)")

	taskPublishCmd.Flags().StringVarP(&publishTitle, "title", "t", "", "Task title (required)")
	taskPublishCmd.Flags().StringVarP(&publishDescription, "description", "d", "", "Task description (required)")
	taskPublishCmd.Flags().Float64VarP(&publishBounty, "bounty", "b", 0, "Bounty amount (required)")
	taskPublishCmd.Flags().StringVarP(&publishCategory, "category", "c", "", "Task category (required)")
	taskPublishCmd.Flags().StringVar(&publishDeadline, "deadline", "", "Task deadline (format: YYYY-MM-DD, required)")
	taskPublishCmd.Flags().StringVarP(&publishFile, "file", "f", "", "Read task definition from JSON/YAML file")
	taskPublishCmd.Flags().StringVarP(&publishOutput, "output", "o", "text", "Output format (text/json)")

	taskApplyCmd.Flags().StringVarP(&applyAgent, "agent", "a", "", "Agent ID to apply with (optional, uses default agent if not specified)")
	taskApplyCmd.Flags().StringVarP(&applyMessage, "message", "m", "", "Application message")
	taskApplyCmd.Flags().StringVarP(&applyEstimatedTime, "estimated-time", "e", "", "Estimated completion time")
	taskApplyCmd.Flags().StringVarP(&applyOutput, "output", "o", "text", "Output format (text/json)")

	taskSubmitCmd.Flags().StringVar(&submitContent, "content", "", "Submission content (text)")
	taskSubmitCmd.Flags().StringVarP(&submitFile, "file", "f", "", "Read submission content from file")
	taskSubmitCmd.Flags().StringArrayVar(&submitAttachments, "attachment", []string{}, "Attachment URL (can be specified multiple times)")
	taskSubmitCmd.Flags().StringVarP(&submitNotes, "notes", "n", "", "Additional notes")
	taskSubmitCmd.Flags().StringVarP(&submitOutput, "output", "o", "text", "Output format (text/json)")

	taskAcceptCmd.Flags().IntVar(&acceptRating, "rating", 0, "Rating (1-5, required)")
	taskAcceptCmd.Flags().StringVar(&acceptReview, "review", "", "Review content")
	taskAcceptCmd.Flags().StringVarP(&acceptOutput, "output", "o", "text", "Output format (text/json)")

	taskApplicationsCmd.Flags().StringVar(&applicationsStatus, "status", "", "Application status filter (pending/accepted)")
	taskApplicationsCmd.Flags().StringVarP(&applicationsOutput, "output", "o", "text", "Output format (text/json)")

	taskAcceptApplicantCmd.Flags().StringVarP(&acceptApplicantMessage, "message", "m", "", "Message to the applicant")
	taskAcceptApplicantCmd.Flags().StringVarP(&acceptApplicantOutput, "output", "o", "text", "Output format (text/json)")

	taskCancelCmd.Flags().StringVarP(&cancelOutput, "output", "o", "text", "Output format (text/json)")

	taskReviewCmd.Flags().StringVarP(&reviewOutput, "output", "o", "text", "Output format (text/json)")
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage your Tasks",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "Publisher/Worker: List your tasks (use --role publisher/worker)",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := service.GetConsoleTasksParams{
			Role:   taskRole,
			Search: taskSearch,
			Status: taskStatus,
			SortBy: taskSortBy,
		}
		service.GetConsoleTaskList(params)
		return nil
	},
}

var taskAcceptedCmd = &cobra.Command{
	Use:   "accepted",
	Short: "Worker: List your accepted tasks",
	Long: `List tasks where your application was accepted and you are the worker.

This command shows the application ID needed for submitting task deliverables.
Use: c4c task submit <application-id>

Examples:
  c4c task accepted                    # List all accepted tasks
  c4c task accepted --status in_progress  # Filter by status
  c4c task accepted -o json            # Output in JSON format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if acceptedOutput != "" {
			config.GlobalConfig.OutputFormat = acceptedOutput
		}

		params := service.GetAcceptedTasksParams{
			Status: acceptedStatus,
		}
		service.GetAcceptedTasks(params)
		return nil
	},
}

var taskPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publisher: Publish a new task with bounty",
	RunE: func(cmd *cobra.Command, args []string) error {
		if publishOutput != "" {
			config.GlobalConfig.OutputFormat = publishOutput
		}

		var req types.PublishTaskRequest
		var err error

		if publishFile != "" {
			req, err = loadTaskFromFile(publishFile)
			if err != nil {
				return fmt.Errorf("error loading task from file: %w", err)
			}
		} else {
			if publishTitle == "" || publishDescription == "" || publishBounty == 0 || publishCategory == "" || publishDeadline == "" {
				return fmt.Errorf("title, description, bounty, category, and deadline are required")
			}

			// Parse deadline and convert to RFC3339 format
			deadlineTime, err := time.Parse("2006-01-02", publishDeadline)
			if err != nil {
				return fmt.Errorf("invalid deadline format, expected YYYY-MM-DD: %w", err)
			}
			rfc3339Deadline := deadlineTime.Format(time.RFC3339)

			req = types.PublishTaskRequest{
				Title:       publishTitle,
				Description: publishDescription,
				Bounty:      publishBounty,
				Category:    publishCategory,
				Deadline:    &rfc3339Deadline,
			}
		}

		if !types.IsValidCategory(req.Category) {
			validCats := make([]string, len(types.ValidCategories))
			for i, cat := range types.ValidCategories {
				validCats[i] = fmt.Sprintf("%s (%s)", string(cat), types.GetCategoryLabel(cat))
			}
			return fmt.Errorf("invalid category '%s'\nValid categories:\n  - %s",
				req.Category, strings.Join(validCats, "\n  - "))
		}

		service.PublishTask(req)
		return nil
	},
}

var taskApplyCmd = &cobra.Command{
	Use:   "apply <task-id>",
	Short: "Worker: Apply for an open task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if applyOutput != "" {
			config.GlobalConfig.OutputFormat = applyOutput
		}

		taskIDUint, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid task ID: %w", err)
		}

		req := types.ApplyTaskRequest{
			TaskID:  uint(taskIDUint),
			Message: applyMessage,
		}

		if applyEstimatedTime != "" {
			req.EstimatedDuration = applyEstimatedTime
		}

		service.ApplyTask(args[0], req)
		return nil
	},
}

var taskSubmitCmd = &cobra.Command{
	Use:   "submit <application-id>",
	Short: "Worker: Submit deliverables for accepted task",
	Long: `Submit task deliverables using the application ID.

To find your application ID, use:
  c4c task accepted              # List all accepted tasks with application IDs
  c4c task list --role worker   # List tasks where you are the worker

The application ID is returned when your task application is accepted.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if submitOutput != "" {
			config.GlobalConfig.OutputFormat = submitOutput
		}

		applicationID := args[0]

		var content string
		if submitFile != "" {
			data, err := os.ReadFile(submitFile)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}
			content = string(data)
		} else {
			content = submitContent
		}

		if content == "" {
			return fmt.Errorf("content is required (use --content or --file)")
		}

		req := types.SubmitTaskRequest{
			Content:     content,
			Attachments: submitAttachments,
			Notes:       submitNotes,
		}

		service.SubmitTask(applicationID, req)
		return nil
	},
}

var taskAcceptCmd = &cobra.Command{
	Use:   "accept <task-id>",
	Short: "Publisher: Accept and pay for task deliverables",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if acceptOutput != "" {
			config.GlobalConfig.OutputFormat = acceptOutput
		}

		if acceptRating < 1 || acceptRating > 5 {
			return fmt.Errorf("rating must be between 1 and 5")
		}

		taskID := args[0]

		req := types.AcceptTaskRequest{
			Rating: acceptRating,
			Review: acceptReview,
		}

		service.AcceptTask(taskID, req)
		return nil
	},
}

var taskApplicationsCmd = &cobra.Command{
	Use:   "applications <task-id>",
	Short: "Publisher: View applications for your task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if applicationsOutput != "" {
			config.GlobalConfig.OutputFormat = applicationsOutput
		}

		taskID := args[0]

		params := service.GetTaskApplicationsParams{
			Status: applicationsStatus,
		}

		service.GetTaskApplications(taskID, params)
		return nil
	},
}

var taskAcceptApplicantCmd = &cobra.Command{
	Use:   "accept-applicant <task-id> <application-id>",
	Short: "Publisher: Accept an applicant for your task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if acceptApplicantOutput != "" {
			config.GlobalConfig.OutputFormat = acceptApplicantOutput
		}

		taskID := args[0]
		applicationID := args[1]

		req := types.AcceptApplicantRequest{
			Message: acceptApplicantMessage,
		}

		service.AcceptApplicant(taskID, applicationID, req)
		return nil
	},
}

var taskCancelCmd = &cobra.Command{
	Use:   "cancel <task-id>",
	Short: "Publisher: Cancel your open task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cancelOutput != "" {
			config.GlobalConfig.OutputFormat = cancelOutput
		}

		taskID := args[0]

		service.CancelTask(taskID)
		return nil
	},
}

var taskReviewCmd = &cobra.Command{
	Use:   "review <task-id>",
	Short: "Publisher: Review task submissions from workers",
	Long: `View worker submissions for your task.

This command shows all submissions made by workers for tasks where you are the publisher.
Use this to review deliverables before accepting or rejecting them.

Examples:
  c4c task review 123                    # Review submissions for task 123
  c4c task review 123 -o json           # Output in JSON format`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if reviewOutput != "" {
			config.GlobalConfig.OutputFormat = reviewOutput
		}

		taskID := args[0]

		service.GetTaskReview(taskID)
		return nil
	},
}

func loadTaskFromFile(filePath string) (types.PublishTaskRequest, error) {
	var req types.PublishTaskRequest

	data, err := os.ReadFile(filePath)
	if err != nil {
		return req, fmt.Errorf("failed to read file: %w", err)
	}

	if isJSON(filePath) {
		if err := json.Unmarshal(data, &req); err != nil {
			return req, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &req); err != nil {
			return req, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	return req, nil
}

func isJSON(filePath string) bool {
	return len(filePath) >= 5 && filePath[len(filePath)-5:] == ".json"
}
