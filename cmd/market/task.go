package market

import (
	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/spf13/cobra"
)

var (
	taskListPage     int
	taskListLimit    int
	taskListSearch   string
	taskListCategory string
	taskListStatus   string
	taskOutputFormat string
)

func init() {
	MarketCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskSearchCmd)

	taskListCmd.Flags().IntVarP(&taskListPage, "page", "p", 1, "Page number")
	taskListCmd.Flags().IntVarP(&taskListLimit, "limit", "l", 9, "Number of items per page")
	taskListCmd.Flags().StringVarP(&taskListSearch, "search", "s", "", "Search keyword")
	taskListCmd.Flags().StringVarP(&taskListCategory, "category", "c", "", "Category filter")
	taskListCmd.Flags().StringVar(&taskListStatus, "status", "", "Status filter (open, in_progress, completed, etc.)")
	taskListCmd.Flags().StringVarP(&taskOutputFormat, "output", "o", "text", "Output format (text/json)")

	taskShowCmd.Flags().StringVarP(&taskOutputFormat, "output", "o", "text", "Output format (text/json)")

	taskSearchCmd.Flags().IntVarP(&taskListPage, "page", "p", 1, "Page number")
	taskSearchCmd.Flags().IntVarP(&taskListLimit, "limit", "l", 9, "Number of items per page")
	taskSearchCmd.Flags().StringVarP(&taskListCategory, "category", "c", "", "Category filter")
	taskSearchCmd.Flags().StringVar(&taskListStatus, "status", "", "Status filter (open, in_progress, completed, etc.)")
	taskSearchCmd.Flags().StringVarP(&taskOutputFormat, "output", "o", "text", "Output format (text/json)")
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Explore Tasks in the market",
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "Worker: List all public tasks in the market",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = taskOutputFormat

		params := service.GetMarketTasksParams{
			Page:     taskListPage,
			Limit:    taskListLimit,
			Search:   taskListSearch,
			Category: taskListCategory,
			Status:   taskListStatus,
		}
		service.GetMarketTaskList(params)
	},
}

var taskShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Worker: Show details of a specific task in the market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = taskOutputFormat
		service.GetMarketTaskDetail(args[0])
	},
}

var taskSearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Worker: Search tasks in the market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = taskOutputFormat

		params := service.GetMarketTasksParams{
			Page:     taskListPage,
			Limit:    taskListLimit,
			Search:   args[0],
			Category: taskListCategory,
			Status:   taskListStatus,
		}
		service.GetMarketTaskList(params)
	},
}
