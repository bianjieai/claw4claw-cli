package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintAgentStatus(w io.Writer, data *types.AgentStatus) {
	if config.GlobalConfig.OutputFormat == "json" {
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintln(w, "{}")
			return
		}
		fmt.Fprintln(w, string(bytes))
		return
	}

	fmt.Fprintln(w, "=== Agent Status ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "DID\t%s\n", data.DID)
	fmt.Fprintf(tw, "Name\t%s\n", data.Name)
	fmt.Fprintf(tw, "Status\t%s\n", data.Status)
	fmt.Fprintf(tw, "Balance\t%s\n", data.Balance)
	fmt.Fprintf(tw, "Reputation Level\t%s\n", data.Reputation.Level)
	fmt.Fprintf(tw, "Reputation Score\t%.1f\n", data.Reputation.Score)
	tw.Flush()
}

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
	fmt.Fprintf(tw, "DID\t%s\n", agent.DID)
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
