package market

import (
	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/spf13/cobra"
)

var (
	serviceListPage     int
	serviceListLimit    int
	serviceListSearch   string
	serviceListCategory string
	serviceListStatus   string
	serviceOutputFormat string
)

func init() {
	MarketCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceShowCmd)
	serviceCmd.AddCommand(serviceSearchCmd)

	serviceListCmd.Flags().IntVarP(&serviceListPage, "page", "p", 1, "Page number")
	serviceListCmd.Flags().IntVarP(&serviceListLimit, "limit", "l", 9, "Number of items per page")
	serviceListCmd.Flags().StringVarP(&serviceListSearch, "search", "s", "", "Search keyword")
	serviceListCmd.Flags().StringVarP(&serviceListCategory, "category", "c", "", "Category filter")
	serviceListCmd.Flags().StringVar(&serviceListStatus, "status", "", "Status filter (active, etc.)")
	serviceListCmd.Flags().StringVarP(&serviceOutputFormat, "output", "o", "text", "Output format (text/json)")

	serviceShowCmd.Flags().StringVarP(&serviceOutputFormat, "output", "o", "text", "Output format (text/json)")

	serviceSearchCmd.Flags().IntVarP(&serviceListPage, "page", "p", 1, "Page number")
	serviceSearchCmd.Flags().IntVarP(&serviceListLimit, "limit", "l", 9, "Number of items per page")
	serviceSearchCmd.Flags().StringVarP(&serviceListCategory, "category", "c", "", "Category filter")
	serviceSearchCmd.Flags().StringVar(&serviceListStatus, "status", "", "Status filter (active, etc.)")
	serviceSearchCmd.Flags().StringVarP(&serviceOutputFormat, "output", "o", "text", "Output format (text/json)")
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Explore Services in the market",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Caller: List all public services in the market",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = serviceOutputFormat

		params := service.GetMarketServicesParams{
			Page:     serviceListPage,
			Limit:    serviceListLimit,
			Search:   serviceListSearch,
			Category: serviceListCategory,
			Status:   serviceListStatus,
		}
		service.GetMarketServiceList(params)
	},
}

var serviceShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Caller: Show details of a specific service in the market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = serviceOutputFormat
		service.GetMarketServiceDetail(args[0])
	},
}

var serviceSearchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Caller: Search services in the market",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.OutputFormat = serviceOutputFormat

		params := service.GetMarketServicesParams{
			Page:     serviceListPage,
			Limit:    serviceListLimit,
			Search:   args[0],
			Category: serviceListCategory,
			Status:   serviceListStatus,
		}
		service.GetMarketServiceList(params)
	},
}
