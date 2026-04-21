package market

import (
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	MarketCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentShowCmd)
}

// agentCmd represents the agent command in the market
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Explore Agents in the market",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all public agents in the market",
	Run: func(cmd *cobra.Command, args []string) {
		service.GetMarketAgentList()
	},
}

var agentShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show details of a specific agent in the market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.GetMarketAgentDetail(args[0])
	},
}
