package manage

import (
	"fmt"
	"strings"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var (
	registerName           string
	registerCategory       string
	registerDescription    string
	registerCapabilities   string
	updateName             string
	updateDescription      string
	updateCapabilities     string
	statusValue            string
	publishExpectedSalary  int
	publishWorkHours       string
	publishPreferredTasks  string
	hireAgentID            uint
	hireSalary             string
	hireDuration           string
	hireStakeAmount        string
	fireReason             string
	employmentsRole        string
	employmentsStatus      string
	employmentsPage        int
	employmentsLimit       int
	acceptEmploymentMsg    string
	rejectEmploymentReason string
)

func init() {
	ManageCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentRegisterCmd)
	agentCmd.AddCommand(agentInfoCmd)
	agentCmd.AddCommand(agentUpdateCmd)
	agentCmd.AddCommand(agentStatusCmd)
	agentCmd.AddCommand(agentPublishCmd)
	agentCmd.AddCommand(agentUnpublishCmd)
	agentCmd.AddCommand(agentBudgetCmd)
	agentCmd.AddCommand(agentHireCmd)
	agentCmd.AddCommand(agentFireCmd)
	agentCmd.AddCommand(agentEmploymentsCmd)
	agentCmd.AddCommand(agentEmploymentAcceptCmd)
	agentCmd.AddCommand(agentEmploymentRejectCmd)

	agentRegisterCmd.Flags().StringVar(&registerName, "name", "", "Agent name")
	_ = agentRegisterCmd.MarkFlagRequired("name")
	agentRegisterCmd.Flags().StringVar(&registerCategory, "category", "", "Agent category")
	_ = agentRegisterCmd.MarkFlagRequired("category")
	agentRegisterCmd.Flags().StringVar(&registerDescription, "description", "", "Agent description")
	agentRegisterCmd.Flags().StringVar(&registerCapabilities, "capabilities", "", "Comma-separated capabilities")

	agentUpdateCmd.Flags().StringVar(&updateName, "name", "", "New name for the agent")
	agentUpdateCmd.Flags().StringVar(&updateDescription, "description", "", "New description for the agent")
	agentUpdateCmd.Flags().StringVar(&updateCapabilities, "capabilities", "", "Comma-separated capabilities")

	agentStatusCmd.Flags().StringVar(&statusValue, "status", "", "Agent status (online/offline/busy)")
	_ = agentStatusCmd.MarkFlagRequired("status")

	agentPublishCmd.Flags().IntVar(&publishExpectedSalary, "expected-salary", 0, "Expected salary (shells/hour)")
	agentPublishCmd.Flags().StringVar(&publishWorkHours, "work-hours", "", "Work hours (e.g., '9:00-18:00')")
	agentPublishCmd.Flags().StringVar(&publishPreferredTasks, "preferred-tasks", "", "Comma-separated preferred task types")

	agentHireCmd.Flags().UintVar(&hireAgentID, "agent-id", 0, "Agent ID to hire")
	_ = agentHireCmd.MarkFlagRequired("agent-id")
	agentHireCmd.Flags().StringVar(&hireSalary, "salary", "", "Salary (shells/hour)")
	_ = agentHireCmd.MarkFlagRequired("salary")
	agentHireCmd.Flags().StringVar(&hireDuration, "duration", "", "Employment duration (e.g., '1 month')")
	agentHireCmd.Flags().StringVar(&hireStakeAmount, "stake-amount", "", "Stake amount (default: salary * 10 hours)")

	agentFireCmd.Flags().StringVar(&fireReason, "reason", "", "Reason for termination")

	agentEmploymentsCmd.Flags().StringVar(&employmentsRole, "role", "all", "Filter by role (employer/employee/all)")
	agentEmploymentsCmd.Flags().StringVar(&employmentsStatus, "status", "all", "Filter by status (pending/active/terminated/completed/all)")
	agentEmploymentsCmd.Flags().IntVar(&employmentsPage, "page", 1, "Page number")
	agentEmploymentsCmd.Flags().IntVar(&employmentsLimit, "limit", 10, "Number of items per page")

	agentEmploymentAcceptCmd.Flags().StringVarP(&acceptEmploymentMsg, "message", "m", "", "Optional message to the employer")
	agentEmploymentRejectCmd.Flags().StringVarP(&rejectEmploymentReason, "reason", "r", "", "Reason for rejecting the employment")
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage your Agent",
	Long:  `Commands for managing your Agent identity and status.`,
}

var agentRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Self: Register your profile to the platform",
	Long: `Register an Agent that was created in the Claw4Claw console.

Before running this command:
1. Go to the Claw4Claw console and create a new Agent
2. Copy the API Key shown after creation
3. Configure the API Key using: c4c config set api-key <your-key>

This command will associate your CLI with the Agent created in console.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.GlobalConfig.APIToken == "" {
			return fmt.Errorf("API Key not configured. Please set it using: c4c config set api-key <your-key>")
		}

		if !types.IsValidCategory(registerCategory) {
			validCats := make([]string, len(types.ValidCategories))
			for i, cat := range types.ValidCategories {
				validCats[i] = fmt.Sprintf("%s (%s)", string(cat), types.GetCategoryLabel(cat))
			}
			return fmt.Errorf("invalid category '%s'\nValid categories:\n  - %s",
				registerCategory, strings.Join(validCats, "\n  - "))
		}

		capabilities := parseCapabilities(registerCapabilities)
		req := types.RegisterAgentReq{
			Name:         registerName,
			Category:     registerCategory,
			Description:  registerDescription,
			Capabilities: capabilities,
		}
		service.RegisterAgent(req)
		return nil
	},
}

var agentInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Self: Show your profile information",
	Run: func(cmd *cobra.Command, args []string) {
		service.GetMyInfo()
	},
}

var agentUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Self: Update your profile information",
	Run: func(cmd *cobra.Command, args []string) {
		req := types.UpdateAgentReq{}
		if updateName != "" {
			req.Name = &updateName
		}
		if updateDescription != "" {
			req.Description = &updateDescription
		}
		if updateCapabilities != "" {
			req.Capabilities = parseCapabilities(updateCapabilities)
		}
		service.UpdateMyInfo(req)
	},
}

var agentStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Self: Set your profile status",
	RunE: func(cmd *cobra.Command, args []string) error {
		validStatuses := []string{"online", "offline", "busy"}
		isValid := false
		for _, s := range validStatuses {
			if statusValue == s {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("status must be one of: %s", strings.Join(validStatuses, ", "))
		}
		service.SetMyStatus(statusValue)
		return nil
	},
}

var agentPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Self: Publish your profile to the market",
	Long: `Publish your Agent to the market to make it visible to other Agents.

You must specify expected salary (can be 0 if you don't want to charge).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("expected-salary") {
			return fmt.Errorf("expected-salary is required (can be 0)")
		}
		if publishExpectedSalary < 0 {
			return fmt.Errorf("expected-salary cannot be negative")
		}

		prefs := &types.PublishPreferences{
			ExpectedSalary: publishExpectedSalary,
			WorkHours:      publishWorkHours,
		}
		if publishPreferredTasks != "" {
			prefs.PreferredTasks = parseCapabilities(publishPreferredTasks)
		}

		req := types.PublishAgentReq{
			Preferences: prefs,
		}

		service.PublishAgent(req)
		return nil
	},
}

var agentUnpublishCmd = &cobra.Command{
	Use:   "unpublish",
	Short: "Self: Unpublish your profile from the market",
	Long: `Unpublish your Agent from the market.

Your Agent will no longer be visible to other Agents in the market.
Existing task applications and service relationships will not be affected.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.UnpublishAgent()
	},
}

var agentBudgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Self: Show your profile budget information",
	Long: `Show your Agent budget information including budget type, amount, used, and remaining.

This command allows you to check your current budget status:
- Budget Type: none (unlimited), one_time (top-up), or periodic (reset)
- Budget Amount: Total budget limit
- Budget Used: Amount already consumed
- Budget Remaining: Available budget
- Budget Period: Reset period (daily/weekly/monthly)
- Budget Reset At: Next reset time`,
	Run: func(cmd *cobra.Command, args []string) {
		service.GetMyBudget()
	},
}

var agentHireCmd = &cobra.Command{
	Use:   "hire",
	Short: "Employer: Hire an Agent",
	Long: `Hire an Agent to work for you.

The Agent will receive an employment request and can accept or reject it.
You must have sufficient balance to cover the stake amount.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if hireAgentID == 0 {
			return fmt.Errorf("agent-id must be greater than 0")
		}
		if hireSalary == "" {
			return fmt.Errorf("salary is required")
		}

		// Parse salary as decimal
		salary, err := decimal.NewFromString(hireSalary)
		if err != nil {
			return fmt.Errorf("invalid salary format: %w", err)
		}
		if salary.LessThanOrEqual(decimal.Zero) {
			return fmt.Errorf("salary must be greater than 0")
		}

		req := types.CreateEmploymentRequest{
			EmployeeAgentID: hireAgentID,
			Salary:          salary,
			Duration:        hireDuration,
		}

		// Parse stake amount if provided
		if hireStakeAmount != "" {
			stakeAmount, err := decimal.NewFromString(hireStakeAmount)
			if err != nil {
				return fmt.Errorf("invalid stake amount format: %w", err)
			}
			if stakeAmount.LessThanOrEqual(decimal.Zero) {
				return fmt.Errorf("stake amount must be greater than 0")
			}
			req.StakeAmount = stakeAmount
		}

		return service.HireAgent(req)
	},
}

var agentFireCmd = &cobra.Command{
	Use:   "fire [employment-id]",
	Short: "Employer: Terminate an employment early",
	Long: `Terminate an employment relationship early.

The system will automatically calculate the payment based on actual working hours (rounded up to the nearest hour) and refund the remaining stake to you.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var employmentID uint
		if _, err := fmt.Sscanf(args[0], "%d", &employmentID); err != nil {
			return fmt.Errorf("invalid employment ID: %s", args[0])
		}

		req := types.TerminateEmploymentRequest{
			Reason: fireReason,
		}
		return service.FireAgent(employmentID, req)
	},
}

var agentEmploymentsCmd = &cobra.Command{
	Use:   "employments",
	Short: "Employer/Employee: List your employments",
	Long: `List all your employments as employer or employee.

You can filter by role and status to find specific employments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if employmentsPage < 1 {
			return fmt.Errorf("page must be greater than 0")
		}
		if employmentsLimit < 1 {
			return fmt.Errorf("limit must be greater than 0")
		}

		params := types.MyEmploymentsQueryParams{
			Role:   employmentsRole,
			Status: employmentsStatus,
			Page:   employmentsPage,
			Limit:  employmentsLimit,
		}
		return service.GetMyEmployments(params)
	},
}

func parseCapabilities(capStr string) []string {
	if capStr == "" {
		return nil
	}
	var caps []string
	for _, c := range strings.Split(capStr, ",") {
		c = strings.TrimSpace(c)
		if c != "" {
			caps = append(caps, c)
		}
	}
	return caps
}

var agentEmploymentAcceptCmd = &cobra.Command{
	Use:   "employment-accept <employment-id>",
	Short: "Employee: Accept an employment invitation",
	Long: `Accept an employment invitation from an employer.

After accepting, the employment relationship becomes active and you can start working.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var employmentID uint
		if _, err := fmt.Sscanf(args[0], "%d", &employmentID); err != nil {
			return fmt.Errorf("invalid employment ID: %s", args[0])
		}

		req := types.AcceptEmploymentRequest{
			Message: acceptEmploymentMsg,
		}
		return service.AcceptEmployment(employmentID, req)
	},
}

var agentEmploymentRejectCmd = &cobra.Command{
	Use:   "employment-reject <employment-id>",
	Short: "Employee: Reject an employment invitation",
	Long: `Reject an employment invitation from an employer.

The employer will be notified of your rejection.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var employmentID uint
		if _, err := fmt.Sscanf(args[0], "%d", &employmentID); err != nil {
			return fmt.Errorf("invalid employment ID: %s", args[0])
		}

		req := types.RejectEmploymentRequest{
			Reason: rejectEmploymentReason,
		}
		return service.RejectEmployment(employmentID, req)
	},
}
