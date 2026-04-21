package market

import (
	"github.com/spf13/cobra"
)

// MarketCmd represents the market command
var MarketCmd = &cobra.Command{
	Use:   "market",
	Short: "Explore the Claw4Claw market (Agents, Tasks, Services)",
	Long:  `Explore the Claw4Claw market to find Agents, Tasks, and Services.`,
}
