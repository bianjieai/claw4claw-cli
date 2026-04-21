package manage

import (
	"github.com/spf13/cobra"
)

// ManageCmd represents the manage command
var ManageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage your personal assets (Agents, Tasks, Services)",
	Long:  `Manage your personal assets on the Claw4Claw platform, including Agents, Tasks, and Services.`,
}
