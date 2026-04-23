package printer

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintMarketServiceList(w io.Writer, list *types.MarketServiceList) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	header := fmt.Sprintf("Market Services (Page %d/%d, Total: %d)", list.Page, list.TotalPages, list.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTitle\tCategory\tPrice\tCompleted\tRating\n")
	for _, s := range list.Data {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\t%.1f\n", s.ID, truncate(s.Title, 30), s.Category, s.Price, s.Completed, s.Rating)
	}
	tw.Flush()
}

func PrintMarketServiceDetail(w io.Writer, svc *types.MarketService) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, svc)
		return
	}

	header := fmt.Sprintf("Service Details (#%d)", svc.ID)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", svc.ID)
	fmt.Fprintf(tw, "Title\t%s\n", svc.Title)
	fmt.Fprintf(tw, "Description\t%s\n", svc.Description)
	fmt.Fprintf(tw, "Category\t%s\n", svc.Category)
	fmt.Fprintf(tw, "Price\t%s\n", svc.Price)
	fmt.Fprintf(tw, "Rating\t%.1f\n", svc.Rating)
	fmt.Fprintf(tw, "Completed\t%d\n", svc.Completed)
	fmt.Fprintf(tw, "Provider\t%s (ID: %d)\n", svc.Provider, svc.ProviderID)
	fmt.Fprintf(tw, "Avg Response Time\t%s\n", svc.AvgResponseTime)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(svc.CreatedAt))

	if svc.InputSchema != "" {
		fmt.Fprintf(tw, "Input Schema\t%s\n", truncate(svc.InputSchema, 60))
	}
	if svc.OutputSchema != "" {
		fmt.Fprintf(tw, "Output Schema\t%s\n", truncate(svc.OutputSchema, 60))
	}

	tw.Flush()
}

func PrintConsoleServiceList(w io.Writer, list *types.ConsoleServiceList) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, list)
		return
	}

	header := fmt.Sprintf("My Published Services (Total: %d)", list.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTitle\tCategory\tPrice\tTotal Calls\tTotal Earnings\tStatus\n")
	for _, s := range list.Data {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", s.ID, truncate(s.Title, 30), s.Category, s.Price, s.TotalCalls, s.TotalEarnings, s.Status)
	}
	tw.Flush()
}

func PrintConsoleServiceDetail(w io.Writer, detail *types.ConsoleServiceDetails) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, detail)
		return
	}

	header := fmt.Sprintf("Service Details (#%s)", detail.ID)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%s\n", detail.ID)
	fmt.Fprintf(tw, "Title\t%s\n", detail.Title)
	fmt.Fprintf(tw, "Provider Agent\t%s (ID: %s)\n", detail.ProviderAgentName, detail.ProviderAgentID)
	fmt.Fprintf(tw, "Price\t%s\n", detail.Price)
	fmt.Fprintf(tw, "Total Calls\t%d\n", detail.TotalCalls)
	fmt.Fprintf(tw, "Total Earnings\t%s\n", detail.TotalEarnings)
	fmt.Fprintf(tw, "Avg Response Time\t%s\n", detail.AvgResponseTime)
	fmt.Fprintf(tw, "Status\t%s\n", detail.Status)
	fmt.Fprintf(tw, "Error Rate\t%s\n", detail.ErrorRate)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(detail.CreatedAt))
	tw.Flush()

	if len(detail.RecentLogs) > 0 {
		fmt.Fprintln(w, "\n=== Recent Logs ===")
		twLogs := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(twLogs, "Log ID\tTimestamp\tStatus\tDuration\tCost\n")
		for _, log := range detail.RecentLogs {
			fmt.Fprintf(twLogs, "%s\t%s\t%s\t%s\t%s\n", log.ID, formatDateTime(log.Timestamp), log.Status, log.Duration, log.Cost)
		}
		twLogs.Flush()
	}
}

func PrintPublishServiceSuccess(w io.Writer, resp *types.PublishServiceResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Published Successfully ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", resp.ID)
	fmt.Fprintf(tw, "Title\t%s\n", resp.Title)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(resp.CreatedAt))
	tw.Flush()
}

func PrintMyServicesList(w io.Writer, resp *types.MyServicesListResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	header := fmt.Sprintf("My Services (Total: %d)", resp.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTitle\tStatus\tPrice\tCompleted\tRating\n")
	for _, s := range resp.Services {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%d\t%.1f\n", s.ID, truncate(s.Title, 30), s.Status, s.Price, s.CompletedCount, s.Rating)
	}
	tw.Flush()
}

func PrintUpdateServiceSuccess(w io.Writer, resp *types.UpdateServiceResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Updated Successfully ===")
	fmt.Fprintf(w, "Message: %s\n", resp.Message)
}

func PrintActivateServiceSuccess(w io.Writer, resp *types.ActivateServiceResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Activated Successfully ===")
	fmt.Fprintf(w, "Message: %s\n", resp.Message)
	fmt.Fprintf(w, "Status: %s\n", resp.Status)
}

func PrintDeactivateServiceSuccess(w io.Writer, resp *types.DeactivateServiceResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Deactivated Successfully ===")
	fmt.Fprintf(w, "Message: %s\n", resp.Message)
	fmt.Fprintf(w, "Status: %s\n", resp.Status)
}

func PrintServiceInvocationList(w io.Writer, resp *types.ServiceInvocationListResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	header := fmt.Sprintf("Service Invocations (Page %d, Total: %d)", resp.Page, resp.Total)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tService\tRole\tStatus\tPrice\tCreated At\n")
	for _, inv := range resp.Invocations {
		completedAt := "-"
		if inv.CompletedAt != nil {
			completedAt = formatDateTime(*inv.CompletedAt)
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", inv.ID, truncate(inv.ServiceTitle, 20), inv.Role, inv.Status, inv.Price, formatDateTime(inv.CreatedAt))
		_ = completedAt
	}
	tw.Flush()
}

func PrintServiceInvocationDetail(w io.Writer, detail *types.ServiceInvocationDetail) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, detail)
		return
	}

	header := fmt.Sprintf("Service Invocation Details (#%d)", detail.ID)
	fmt.Fprintln(w, "=== "+header+" ===")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", detail.ID)
	fmt.Fprintf(tw, "Service ID\t%d\n", detail.ServiceID)
	fmt.Fprintf(tw, "Service Title\t%s\n", detail.ServiceTitle)
	fmt.Fprintf(tw, "Role\t%s\n", detail.Role)
	fmt.Fprintf(tw, "Status\t%s\n", detail.Status)
	fmt.Fprintf(tw, "Price\t%s\n", detail.Price)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(detail.CreatedAt))
	if detail.CompletedAt != nil {
		fmt.Fprintf(tw, "Completed At\t%s\n", formatDateTime(*detail.CompletedAt))
	}
	if detail.Rating != nil {
		fmt.Fprintf(tw, "Rating\t%d\n", *detail.Rating)
	}
	if detail.Review != "" {
		fmt.Fprintf(tw, "Review\t%s\n", detail.Review)
	}
	if detail.TimeoutAt != nil {
		fmt.Fprintf(tw, "Timeout At\t%s\n", formatDateTime(*detail.TimeoutAt))
	}
	tw.Flush()

	if len(detail.Input) > 0 {
		fmt.Fprintln(w, "\n=== Input ===")
		printJSON(w, detail.Input)
	}

	if len(detail.Output) > 0 {
		fmt.Fprintln(w, "\n=== Output ===")
		printJSON(w, detail.Output)
	}

	if len(detail.Attachments) > 0 {
		fmt.Fprintln(w, "\n=== Attachments ===")
		for _, a := range detail.Attachments {
			fmt.Fprintf(w, "- %s\n", a)
		}
	}
}

func PrintInvokeServiceSuccess(w io.Writer, resp *types.InvokeServiceResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Invoked Successfully ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Invocation ID\t%d\n", resp.InvocationID)
	fmt.Fprintf(tw, "Service ID\t%d\n", resp.ServiceID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Price\t%s\n", resp.Price)
	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTime(resp.CreatedAt))
	tw.Flush()
}

func PrintSubmitServiceResultSuccess(w io.Writer, resp *types.SubmitServiceResultResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Result Submitted Successfully ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Invocation ID\t%d\n", resp.InvocationID)
	fmt.Fprintf(tw, "Status\t%s\n", resp.Status)
	fmt.Fprintf(tw, "Completed At\t%s\n", formatDateTime(resp.CompletedAt))
	tw.Flush()
}

func PrintReviewServiceInvocationSuccess(w io.Writer, resp *types.ReviewServiceInvocationResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "=== Service Invocation Reviewed Successfully ===")
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "Invocation ID\t%d\n", resp.InvocationID)
	fmt.Fprintf(tw, "Message\t%s\n", resp.Message)
	tw.Flush()
}
