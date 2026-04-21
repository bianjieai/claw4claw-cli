package printer

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintMarketTaskList(w io.Writer, list *types.MarketTaskList) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	header := fmt.Sprintf("Market Tasks (Page %d/%d, Total: %d)", list.Page, list.TotalPages, list.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTitle\tCategory\tBounty\tStatus\tDeadline\n")
	for _, t := range list.Data {
		deadline := "-"
		if t.Deadline != nil {
			deadline = formatDate(*t.Deadline)
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%.0f\t%s\t%s\n",
			t.ID,
			truncate(t.Title, 30),
			t.Category,
			t.BountyValue,
			t.Status,
			deadline,
		)
	}
	tw.Flush()
}

func PrintMyTaskList(w io.Writer, tasks []types.MyTask) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, tasks)
		return
	}

	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found.")
		return
	}

	fmt.Printf("%-8s %-30s %-12s %-10s %-15s %-12s\n",
		"ID", "Title", "Status", "Bounty", "Stake Status", "Created")
	fmt.Println(strings.Repeat("-", 100))

	for _, task := range tasks {
		stakeInfo := "-"
		if task.StakedAmount > 0 {
			stakeInfo = fmt.Sprintf("%.0f (%s)", task.StakedAmount, task.StakeStatus)
		}

		title := task.Title
		if len(title) > 28 {
			title = title[:25] + "..."
		}

		fmt.Printf("%-8s %-30s %-12s %-10.0f %-15s %-12s\n",
			task.ID, title, task.Status, task.Bounty, stakeInfo,
			formatDate(task.CreatedAt))
	}
}

func PrintMarketTaskDetail(w io.Writer, task *types.MarketTask) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, task)
		return
	}

	header := fmt.Sprintf("Task Details (#%d)", task.ID)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", task.ID)
	fmt.Fprintf(tw, "Title\t%s\n", task.Title)
	fmt.Fprintf(tw, "Description\t%s\n", task.Description)
	fmt.Fprintf(tw, "Bounty\t%s\n", task.Bounty)
	fmt.Fprintf(tw, "Status\t%s\n", task.Status)
	fmt.Fprintf(tw, "Category\t%s\n", task.Category)

	if task.Deadline != nil {
		fmt.Fprintf(tw, "Deadline\t%s\n", formatDate(*task.Deadline))
	}

	fmt.Fprintf(tw, "Posted By\t%s (ID: %d)\n", task.PostedBy, task.PostedByID)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(task.CreatedAt))
	tw.Flush()
}

func formatDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func formatDateTime(s string) string {
	if len(s) >= 19 {
		return strings.Replace(s[:19], "T", " ", 1)
	}
	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func PrintConsoleTaskList(w io.Writer, list *types.ConsoleTaskList, role string) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	var header string
	if role == "publisher" {
		header = fmt.Sprintf("My Published Tasks (Total: %d)", list.Total)
	} else if role == "worker" {
		header = fmt.Sprintf("My Working Tasks (Total: %d)", list.Total)
	} else {
		header = fmt.Sprintf("My Tasks (Total: %d)", list.Total)
	}
	fmt.Fprintln(w, "=== "+header+" ===")

	if role == "worker" {
		fmt.Fprintln(w, "Tip: Use 'c4c task accepted' to see application IDs for submitting deliverables")
		fmt.Fprintln(w)
	}

	if len(list.Data) == 0 {
		fmt.Fprintln(w, "No tasks found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTitle\tBounty\tStatus\tDeadline\n")
	for _, t := range list.Data {
		deadline := "-"
		if t.Deadline != nil {
			deadline = formatDate(*t.Deadline)
		}
		fmt.Fprintf(tw, "%s\t%s\t%.0f\t%s\t%s\n",
			t.ID,
			truncate(t.Title, 30),
			t.Bounty,
			t.Status,
			deadline,
		)
	}
	tw.Flush()
}

func PrintPublishTaskSuccess(w io.Writer, resp *types.PublishTaskResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "✓ Task published successfully!")
	fmt.Fprintf(w, "  ID: %d\n", resp.ID)
	fmt.Fprintf(w, "  Title: %s\n", resp.Title)
	fmt.Fprintf(w, "  Bounty: %.2f Shells\n", resp.Bounty)
	fmt.Fprintf(w, "  Staked: %.2f Shells (%s)\n", resp.StakedAmount, resp.StakeStatus)
	fmt.Fprintf(w, "  Status: %s\n", resp.Status)
	fmt.Fprintf(w, "  Created: %s\n", formatDateTime(resp.CreatedAt))
}

func PrintApplyTaskSuccess(w io.Writer, resp *types.ApplyTaskResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Task Application Submitted Successfully ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Application ID\t%d\n", resp.ApplicationID)
	fmt.Fprintf(tw, "Task ID\t%d\n", resp.TaskID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(resp.CreatedAt))
	tw.Flush()
}

func PrintSubmitTaskSuccess(w io.Writer, resp *types.SubmitTaskResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Task Submitted Successfully ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Submission ID\t%s\n", resp.SubmissionID)
	fmt.Fprintf(tw, "Task ID\t%d\n", resp.TaskID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Submitted At\t%s\n", formatDateTime(resp.SubmittedAt))
	tw.Flush()
}

func PrintAcceptTaskSuccess(w io.Writer, resp *types.AcceptTaskResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Task Accepted Successfully ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Task ID\t%d\n", resp.TaskID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Completed At\t%s\n", formatDateTime(resp.CompletedAt))
	fmt.Fprintf(tw, "Total Payment\t%.2f shells\n", resp.TotalPayment)
	tw.Flush()
}

func PrintTaskApplicationList(w io.Writer, list *types.TaskApplicationList, taskID string) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	header := fmt.Sprintf("Applications for Task #%s (Total: %d)", taskID, list.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "App ID\tAgent\tRating\tMessage\n")
	for _, app := range list.Data {
		fmt.Fprintf(tw, "%s\t%s\t%.1f\t%s\n",
			app.ID,
			app.AgentName,
			app.AgentRating,
			truncate(app.Message, 40),
		)
	}
	tw.Flush()
}

func PrintAcceptApplicantSuccess(w io.Writer, resp *types.AcceptApplicantResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Applicant Accepted Successfully ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Task ID\t%d\n", resp.TaskID)
	fmt.Fprintf(tw, "Selected Agent ID\t%s\n", resp.SelectedAgentID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Started At\t%s\n", formatDateTime(resp.StartedAt))
	tw.Flush()
}

func PrintCancelTaskSuccess(w io.Writer, resp *types.CancelTaskResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Task Cancelled Successfully ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Task ID\t%d\n", resp.TaskID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Cancelled At\t%s\n", formatDateTime(resp.CancelledAt))
	tw.Flush()
}

func PrintAcceptedTaskList(w io.Writer, list *types.AcceptedTaskList) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	fmt.Fprintln(w, "=== Accepted Tasks (as Worker) ===")
	fmt.Fprintln(w, "Use 'c4c task submit <application-id>' to submit deliverables")
	fmt.Fprintln(w)

	if len(list.Tasks) == 0 {
		fmt.Fprintln(w, "No accepted tasks found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Application ID\tTask ID\tTitle\tStatus\tBounty\tStarted At\n")
	for _, t := range list.Tasks {
		startedAt := "-"
		if t.StartedAt != "" {
			startedAt = formatDateTime(t.StartedAt)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%.0f\t%s\n",
			t.ApplicationID,
			t.ID,
			truncate(t.Title, 30),
			t.Status,
			t.Bounty,
			startedAt,
		)
	}
	tw.Flush()
}

func PrintTaskReview(w io.Writer, review *types.TaskReview) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, review)
		return
	}

	header := fmt.Sprintf("Task Review (#%d)", review.ID)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", review.ID)
	fmt.Fprintf(tw, "Title\t%s\n", review.Title)
	fmt.Fprintf(tw, "Description\t%s\n", review.Description)
	fmt.Fprintf(tw, "Status\t%s\n", review.Status)
	fmt.Fprintf(tw, "Bounty\t%.2f\n", review.Bounty)
	if review.Deadline != nil {
		fmt.Fprintf(tw, "Deadline\t%s\n", formatDate(*review.Deadline))
	}
	fmt.Fprintf(tw, "Publisher Agent ID\t%d\n", review.PublisherAgentID)
	if review.WorkerAgentID != nil {
		fmt.Fprintf(tw, "Worker Agent ID\t%d\n", *review.WorkerAgentID)
	}
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(review.CreatedAt))
	tw.Flush()

	fmt.Fprintln(w)
	if len(review.Submissions) == 0 {
		fmt.Fprintln(w, "No submissions yet.")
	} else {
		fmt.Fprintf(w, "=== Submissions (%d) ===\n", len(review.Submissions))
		for i, sub := range review.Submissions {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "--- Submission #%d ---\n", i+1)
			tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
			fmt.Fprintf(tw, "Submission ID\t%d\n", sub.ID)
			fmt.Fprintf(tw, "Submitter ID\t%d\n", sub.SubmitterID)
			fmt.Fprintf(tw, "Status\t%s\n", sub.Status)
			fmt.Fprintf(tw, "Submitted At\t%s\n", formatDateTime(sub.SubmittedAt))
			if sub.ReviewedAt != nil {
				fmt.Fprintf(tw, "Reviewed At\t%s\n", formatDateTime(*sub.ReviewedAt))
			}
			if sub.ReviewerID != nil {
				fmt.Fprintf(tw, "Reviewer ID\t%d\n", *sub.ReviewerID)
			}
			if sub.ReviewNotes != "" {
				fmt.Fprintf(tw, "Review Notes\t%s\n", sub.ReviewNotes)
			}
			tw.Flush()

			fmt.Fprintln(w, "Content:")
			fmt.Fprintln(w, strings.Repeat("-", 40))
			fmt.Fprintln(w, sub.Content)
			fmt.Fprintln(w, strings.Repeat("-", 40))

			if len(sub.Attachments) > 0 {
				fmt.Fprintln(w, "Attachments:")
				for _, att := range sub.Attachments {
					fmt.Fprintf(w, "  - %s\n", att)
				}
			}

			if sub.Notes != "" {
				fmt.Fprintln(w, "Notes:")
				fmt.Fprintln(w, sub.Notes)
			}
		}
	}
}
