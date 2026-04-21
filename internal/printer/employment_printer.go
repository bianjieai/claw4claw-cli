package printer

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func PrintCreateEmployment(w io.Writer, resp *types.CreateEmploymentResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "Employment created successfully!")
	fmt.Fprintf(w, "Employment ID: %d\n", resp.ID)
	fmt.Fprintf(w, "Status: %s\n", resp.Status)
	fmt.Fprintf(w, "Staked Amount: %.2f shells\n", resp.StakedAmount)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Waiting for agent to accept...")
}

func PrintTerminateEmployment(w io.Writer, resp *types.TerminateEmploymentResponse) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	fmt.Fprintln(w, "Employment terminated successfully!")
	fmt.Fprintf(w, "Employment ID: %d\n", resp.ID)
	fmt.Fprintf(w, "Status: %s\n", resp.Status)
	fmt.Fprintf(w, "Total Duration: %s\n", formatDuration(resp.TotalDuration))
	fmt.Fprintf(w, "Billed Hours: %d\n", resp.BilledHours)
	fmt.Fprintf(w, "Total Payment: %.2f shells\n", resp.TotalPayment)
	fmt.Fprintf(w, "Refund Amount: %.2f shells\n", resp.RefundAmount)
}

func PrintEmploymentList(w io.Writer, resp *types.MyEmploymentsListResponse, role string) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, resp)
		return
	}

	roleTitle := "All"
	if role == "employer" {
		roleTitle = "Employer"
	} else if role == "employee" {
		roleTitle = "Employee"
	}

	fmt.Fprintf(w, "My Employments (as %s):\n", roleTitle)
	fmt.Fprintln(w)

	if len(resp.Employments) == 0 {
		fmt.Fprintln(w, "No employments found.")
		return
	}

	tw := newTabWriter(w)
	defer tw.Flush()

	fmt.Fprintf(tw, "ID\tAgent\tStatus\tSalary\tDuration\tStarted At\n")
	for _, emp := range resp.Employments {
		agentName := emp.EmployeeAgentName
		if role == "employee" {
			agentName = emp.EmployerAgentName
		}
		startTime := "-"
		if emp.StartTime != nil {
			startTime = formatDateTimeFull(*emp.StartTime)
		}
		duration := emp.Duration
		if duration == "" {
			duration = "-"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%.0f/hr\t%s\t%s\n",
			emp.ID, agentName, emp.Status, emp.Salary, duration, startTime)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "Total: %d employments\n", resp.Total)
}

func PrintEmploymentDetail(w io.Writer, emp *types.EmploymentDetail) {
	if config.GlobalConfig.OutputFormat == "json" {
		printJSON(w, emp)
		return
	}

	fmt.Fprintln(w, "=== Employment Detail ===")
	tw := newTabWriter(w)
	defer tw.Flush()

	fmt.Fprintf(tw, "Property\tValue\n")
	fmt.Fprintf(tw, "ID\t%d\n", emp.ID)
	fmt.Fprintf(tw, "Role\t%s\n", emp.Role)
	fmt.Fprintf(tw, "Status\t%s\n", emp.Status)
	fmt.Fprintf(tw, "Salary\t%.2f shells/hr\n", emp.Salary)
	fmt.Fprintf(tw, "Staked Amount\t%.2f shells\n", emp.StakedAmount)
	fmt.Fprintf(tw, "Stake Status\t%s\n", emp.StakeStatus)

	if emp.Duration != "" {
		fmt.Fprintf(tw, "Duration\t%s\n", emp.Duration)
	}

	fmt.Fprintf(tw, "Total Duration\t%s\n", formatDuration(emp.TotalDuration))

	if emp.StartTime != nil {
		fmt.Fprintf(tw, "Started At\t%s\n", formatDateTimeFull(*emp.StartTime))
	}

	if emp.EndTime != nil {
		fmt.Fprintf(tw, "Ended At\t%s\n", formatDateTimeFull(*emp.EndTime))
	}

	fmt.Fprintf(tw, "Created At\t%s\n", formatDateTimeFull(emp.CreatedAt))

	fmt.Fprintln(w)
	fmt.Fprintln(w, "=== Employer Agent ===")
	fmt.Fprintf(tw, "ID\t%d\n", emp.EmployerAgentID)
	fmt.Fprintf(tw, "Name\t%s\n", emp.EmployerAgentName)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "=== Employee Agent ===")
	fmt.Fprintf(tw, "ID\t%d\n", emp.EmployeeAgentID)
	fmt.Fprintf(tw, "Name\t%s\n", emp.EmployeeAgentName)
}

func formatDuration(seconds int64) string {
	if seconds == 0 {
		return "0s"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	if minutes > 0 {
		if secs > 0 {
			return fmt.Sprintf("%dm %ds", minutes, secs)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%ds", secs)
}

func formatDateTimeFull(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02 15:04:05")
}

func newTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
}
