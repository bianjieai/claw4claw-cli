package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintAgentInfo(w io.Writer, agent *types.AgentInfo) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, agent)
		return
	}

	fmt.Fprintln(w, "=== Agent Info ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", agent.ID)
	fmt.Fprintf(tw, "Name\t%s\n", agent.Name)
	fmt.Fprintf(tw, "Category\t%s\n", agent.Category)
	fmt.Fprintf(tw, "Status\t%s\n", agent.Status)
	fmt.Fprintf(tw, "Market Visibility\t%s\n", agent.MarketVisibility)
	if agent.MarketPublishedAt != nil {
		fmt.Fprintf(tw, "Market Published At\t%s\n", *agent.MarketPublishedAt)
	}
	fmt.Fprintf(tw, "Description\t%s\n", agent.Description)
	fmt.Fprintf(tw, "Staked\t%s\n", agent.Staked)
	fmt.Fprintf(tw, "Earned\t%s\n", agent.Earned)
	fmt.Fprintf(tw, "Rating\t%.1f\n", agent.Rating)
	fmt.Fprintf(tw, "Completed Tasks\t%d\n", agent.CompletedTasks)
	fmt.Fprintf(tw, "Created At\t%s\n", agent.CreatedAt)
	tw.Flush()
}

func PrintMarketAgentList(w io.Writer, list *types.MarketAgentList) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	fmt.Fprintln(w, "=== Market Agents ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tName\tCategory\tPrice\tRating\tStatus\n")
	for _, a := range list.Data {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%.1f\t%s\n", a.ID, a.Name, a.Category, a.Price, a.Rating, a.Status)
	}
	tw.Flush()
}

func PrintMarketAgentDetail(w io.Writer, agent *types.MarketAgent) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, agent)
		return
	}

	fmt.Fprintln(w, "=== Market Agent Details ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", agent.ID)
	fmt.Fprintf(tw, "Name\t%s\n", agent.Name)
	fmt.Fprintf(tw, "Category\t%s\n", agent.Category)
	fmt.Fprintf(tw, "Status\t%s\n", agent.Status)
	fmt.Fprintf(tw, "Price\t%s\n", agent.Price)
	fmt.Fprintf(tw, "Rating\t%.1f\n", agent.Rating)
	fmt.Fprintf(tw, "Description\t%s\n", agent.Description)
	if agent.Staked != "" {
		fmt.Fprintf(tw, "Staked\t%s\n", agent.Staked)
	}
	if agent.Uptime != "" {
		fmt.Fprintf(tw, "Uptime\t%s\n", agent.Uptime)
	}
	tw.Flush()
}

func printJSON(w io.Writer, data any) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintln(w, "{}")
		return
	}
	fmt.Fprintln(w, string(bytes))
}

func PrintBudgetInfo(w io.Writer, budget *types.BudgetInfo) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, budget)
		return
	}

	fmt.Fprintln(w, "=== Agent Budget ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")

	budgetTypeLabel := budget.BudgetType
	switch budget.BudgetType {
	case "none":
		budgetTypeLabel = "无限 (无限制)"
	case "one_time":
		budgetTypeLabel = "不限次充值"
	case "periodic":
		budgetTypeLabel = "定期重置"
	}
	fmt.Fprintf(tw, "Budget Type\t%s\n", budgetTypeLabel)

	if budget.BudgetAmount != nil {
		fmt.Fprintf(tw, "Budget Amount\t%s\n", *budget.BudgetAmount)
	}

	fmt.Fprintf(tw, "Budget Used\t%s\n", budget.BudgetUsed)

	if budget.BudgetRemaining != nil {
		fmt.Fprintf(tw, "Budget Remaining\t%s\n", *budget.BudgetRemaining)
	}

	if budget.BudgetPeriod != "" {
		periodLabel := budget.BudgetPeriod
		switch budget.BudgetPeriod {
		case "daily":
			periodLabel = "每日"
		case "weekly":
			periodLabel = "每周"
		case "monthly":
			periodLabel = "每月"
		}
		fmt.Fprintf(tw, "Budget Period\t%s\n", periodLabel)
	}

	if budget.BudgetResetAt != nil {
		fmt.Fprintf(tw, "Budget Reset At\t%s\n", *budget.BudgetResetAt)
	}

	tw.Flush()
}

func FormatNotificationEvent(notif types.WebSocketNotificationMessage) string {
	switch notif.Domain {
	case "task":
		return formatTaskNotification(notif)
	case "service_invocation":
		return formatServiceNotification(notif)
	case "employment":
		return formatEmploymentNotification(notif)
	default:
		return fmt.Sprintf("Event: %s", notif.Event)
	}
}

func formatTaskNotification(notif types.WebSocketNotificationMessage) string {
	switch notif.Event {
	case "task_application":
		var data struct {
			TaskTitle          string `json:"taskTitle"`
			ApplicantAgentName string `json:"applicantAgentName"`
			Bounty             string `json:"bounty"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' applied for your task '%s' (Bounty: %s)", data.ApplicantAgentName, data.TaskTitle, data.Bounty)
		}
	case "task_application_accepted":
		var data struct {
			TaskTitle          string `json:"taskTitle"`
			PublisherAgentName string `json:"publisherAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Your application for task '%s' was accepted by '%s'", data.TaskTitle, data.PublisherAgentName)
		}
	case "task_application_rejected":
		var data struct {
			TaskTitle          string `json:"taskTitle"`
			PublisherAgentName string `json:"publisherAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Your application for task '%s' was rejected by '%s'", data.TaskTitle, data.PublisherAgentName)
		}
	case "task_submitted":
		var data struct {
			TaskTitle          string `json:"taskTitle"`
			SubmitterAgentName string `json:"submitterAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Worker '%s' submitted deliverables for task '%s'", data.SubmitterAgentName, data.TaskTitle)
		}
	case "task_completed":
		var data struct {
			TaskTitle string `json:"taskTitle"`
			Bounty    string `json:"bounty"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Task '%s' is completed and accepted. Bounty (%s) has been released.", data.TaskTitle, data.Bounty)
		}
	case "task_cancelled":
		var data struct {
			TaskTitle string `json:"taskTitle"`
			Reason    string `json:"reason"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Task '%s' was cancelled. Reason: %s", data.TaskTitle, data.Reason)
		}
	}
	return fmt.Sprintf("Task Event: %s", notif.Event)
}

func formatServiceNotification(notif types.WebSocketNotificationMessage) string {
	switch notif.Event {
	case "service_invoked":
		var data struct {
			ServiceTitle    string `json:"serviceTitle"`
			CallerAgentName string `json:"callerAgentName"`
			Price           string `json:"price"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' invoked your service '%s' (Price: %s)", data.CallerAgentName, data.ServiceTitle, data.Price)
		}
	case "result_submitted":
		var data struct {
			ServiceTitle      string `json:"serviceTitle"`
			ProviderAgentName string `json:"providerAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Provider '%s' submitted result for service '%s'", data.ProviderAgentName, data.ServiceTitle)
		}
	case "invocation_failed":
		var data struct {
			ServiceTitle string `json:"serviceTitle"`
			Error        string `json:"error"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Service invocation '%s' failed. Error: %s", data.ServiceTitle, data.Error)
		}
	case "invocation_timeout":
		var data struct {
			ServiceTitle string `json:"serviceTitle"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Service invocation '%s' timed out", data.ServiceTitle)
		}
	case "invocation_reviewed":
		var data struct {
			ServiceTitle    string `json:"serviceTitle"`
			CallerAgentName string `json:"callerAgentName"`
			Rating          int    `json:"rating"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' reviewed your service '%s' with rating %d", data.CallerAgentName, data.ServiceTitle, data.Rating)
		}
	}
	return fmt.Sprintf("Service Event: %s", notif.Event)
}

func formatEmploymentNotification(notif types.WebSocketNotificationMessage) string {
	switch notif.Event {
	case "employment_offered":
		var data struct {
			EmployerAgentName string `json:"employerAgentName"`
			Salary            string `json:"salary"`
			Duration          string `json:"duration"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' offered you an employment (Salary: %s, Duration: %s)", data.EmployerAgentName, data.Salary, data.Duration)
		}
	case "employment_accepted":
		var data struct {
			EmployeeAgentName string `json:"employeeAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' accepted your employment offer", data.EmployeeAgentName)
		}
	case "employment_rejected":
		var data struct {
			EmployeeAgentName string `json:"employeeAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Agent '%s' rejected your employment offer", data.EmployeeAgentName)
		}
	case "employment_terminated":
		var data struct {
			EmployerAgentName string `json:"employerAgentName"`
			Reason            string `json:"reason"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Your employment with '%s' was terminated. Reason: %s", data.EmployerAgentName, data.Reason)
		}
	case "employment_completed":
		var data struct {
			EmployerAgentName string `json:"employerAgentName"`
		}
		if err := json.Unmarshal(notif.Data, &data); err == nil {
			return fmt.Sprintf("Your employment with '%s' has completed", data.EmployerAgentName)
		}
	}
	return fmt.Sprintf("Employment Event: %s", notif.Event)
}

