package manage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/service"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml.v3"
)

var (
	serviceSearch string
	serviceStatus string
	serviceSortBy string

	publishServiceTitle         string
	publishServiceDescription   string
	publishServiceCategory      string
	publishServicePrice         float64
	publishServiceAvgResponseMs int
	publishServiceFile          string
	publishServiceOutput        string

	updateServiceTitle       string
	updateServiceDescription string
	updateServicePrice       float64
)

func init() {
	ManageCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceShowCmd)
	serviceCmd.AddCommand(servicePublishCmd)
	serviceCmd.AddCommand(serviceUpdateCmd)
	serviceCmd.AddCommand(serviceUnpublishCmd)

	serviceListCmd.Flags().StringVarP(&serviceSearch, "search", "s", "", "Search keyword")
	serviceListCmd.Flags().StringVar(&serviceStatus, "status", "all", "Status filter (all/active/inactive)")
	serviceListCmd.Flags().StringVar(&serviceSortBy, "sort-by", "time", "Sort by (time/price/calls)")

	servicePublishCmd.Flags().StringVarP(&publishServiceTitle, "title", "t", "", "Service title (required if not using file)")
	servicePublishCmd.Flags().StringVarP(&publishServiceDescription, "description", "d", "", "Service description (required if not using file)")
	servicePublishCmd.Flags().StringVarP(&publishServiceCategory, "category", "c", "", "Service category (required if not using file)")
	servicePublishCmd.Flags().Float64VarP(&publishServicePrice, "price", "p", 0, "Service price (required if not using file)")
	servicePublishCmd.Flags().IntVarP(&publishServiceAvgResponseMs, "avg-response-ms", "a", 0, "Average response time in ms (required if not using file)")
	servicePublishCmd.Flags().StringVarP(&publishServiceFile, "file", "f", "", "Read service definition from JSON/YAML file")
	servicePublishCmd.Flags().StringVarP(&publishServiceOutput, "output", "o", "text", "Output format (text/json)")

	serviceUpdateCmd.Flags().StringVarP(&updateServiceTitle, "title", "t", "", "Service title")
	serviceUpdateCmd.Flags().StringVarP(&updateServiceDescription, "description", "d", "", "Service description")
	serviceUpdateCmd.Flags().Float64VarP(&updateServicePrice, "price", "p", 0, "Service price")
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage your Services",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Provider: List your published services",
	Run: func(cmd *cobra.Command, args []string) {
		params := service.GetConsoleServicesParams{
			Search: serviceSearch,
			Status: serviceStatus,
			SortBy: serviceSortBy,
		}
		service.GetConsoleServiceList(params)
	},
}

var serviceShowCmd = &cobra.Command{
	Use:   "show <service-id>",
	Short: "Provider: Show details of your service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service.GetConsoleServiceDetail(args[0])
	},
}

var servicePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Provider: Publish a new service to market",
	RunE: func(cmd *cobra.Command, args []string) error {
		if publishServiceOutput != "" {
			config.GlobalConfig.OutputFormat = publishServiceOutput
		}

		var req types.PublishServiceRequest
		var err error

		if publishServiceFile != "" {
			req, err = loadServiceFromFile(publishServiceFile)
			if err != nil {
				return fmt.Errorf("error loading service from file: %w", err)
			}
		} else {
			if publishServiceTitle == "" || publishServiceDescription == "" || publishServiceCategory == "" || publishServicePrice == 0 || publishServiceAvgResponseMs == 0 {
				return fmt.Errorf("title, description, category, price, and avg-response-ms are required")
			}

			req = types.PublishServiceRequest{
				Title:         publishServiceTitle,
				Description:   publishServiceDescription,
				Category:      publishServiceCategory,
				Price:         publishServicePrice,
				AvgResponseMs: publishServiceAvgResponseMs,
			}
		}

		if !types.IsValidCategory(req.Category) {
			validCats := make([]string, len(types.ValidCategories))
			for i, cat := range types.ValidCategories {
				validCats[i] = fmt.Sprintf("%s (%s)", string(cat), types.GetCategoryLabel(cat))
			}
			return fmt.Errorf("invalid category '%s'\nValid categories:\n  - %s",
				req.Category, strings.Join(validCats, "\n  - "))
		}

		service.PublishService(req)
		return nil
	},
}

var serviceUpdateCmd = &cobra.Command{
	Use:   "update <service-id>",
	Short: "Provider: Update your service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceID, err := parseServiceID(args[0])
		if err != nil {
			return err
		}

		if updateServiceTitle == "" && updateServiceDescription == "" && updateServicePrice == 0 {
			return fmt.Errorf("at least one of --title, --description, or --price must be specified")
		}

		req := types.UpdateServiceRequest{}

		if updateServiceTitle != "" {
			req.Title = &updateServiceTitle
		}
		if updateServiceDescription != "" {
			req.Description = &updateServiceDescription
		}
		if updateServicePrice > 0 {
			req.Price = &updateServicePrice
		}

		service.UpdateService(serviceID, req)
		return nil
	},
}

func parseServiceID(idStr string) (int, error) {
	var result int
	_, err := fmt.Sscanf(idStr, "%d", &result)
	if err != nil {
		return 0, fmt.Errorf("invalid service ID: %s", idStr)
	}
	return result, nil
}

func loadServiceFromFile(filePath string) (types.PublishServiceRequest, error) {
	var req types.PublishServiceRequest

	data, err := os.ReadFile(filePath)
	if err != nil {
		return req, fmt.Errorf("failed to read file: %w", err)
	}

	type tmpReq struct {
		Title         string      `json:"title" yaml:"title"`
		Description   string      `json:"description" yaml:"description"`
		Category      string      `json:"category" yaml:"category"`
		Price         float64     `json:"price" yaml:"price"`
		AvgResponseMs int         `json:"avgResponseMs" yaml:"avgResponseMs"`
		InputSchema   interface{} `json:"inputSchema" yaml:"inputSchema"`
		OutputSchema  interface{} `json:"outputSchema" yaml:"outputSchema"`
	}

	var tmp tmpReq

	if isJSON(filePath) {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return req, fmt.Errorf("failed to parse JSON: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &tmp); err != nil {
			return req, fmt.Errorf("failed to parse YAML: %w", err)
		}
	}

	req.Title = tmp.Title
	req.Description = tmp.Description
	req.Category = tmp.Category
	req.Price = tmp.Price
	req.AvgResponseMs = tmp.AvgResponseMs

	if tmp.InputSchema != nil {
		if inputSchema, ok := tmp.InputSchema.(map[string]interface{}); ok {
			req.InputSchema = inputSchema
		}
	}

	if tmp.OutputSchema != nil {
		if outputSchema, ok := tmp.OutputSchema.(map[string]interface{}); ok {
			req.OutputSchema = outputSchema
		}
	}

	return req, nil
}

var serviceUnpublishCmd = &cobra.Command{
	Use:   "unpublish <service-id>",
	Short: "Provider: Unpublish your service from market",
	Long: `Unpublish (deactivate) a service from the market.

The service will no longer be visible to other agents in the market.
Existing service invocations will not be affected.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceID, err := parseServiceID(args[0])
		if err != nil {
			return err
		}
		service.DeactivateService(serviceID)
		return nil
	},
}
