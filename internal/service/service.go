package service

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/printer"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

type GetConsoleServicesParams struct {
	Search string
	Status string
	SortBy string
}

func GetConsoleServiceList(params GetConsoleServicesParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.ConsoleServicesQueryParams{
		Search: params.Search,
		Status: params.Status,
		SortBy: params.SortBy,
	}

	list, err := apiClient.GetConsoleServices(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching console services: %v\n", err)
		os.Exit(1)
	}
	printer.PrintConsoleServiceList(os.Stdout, list)
}

func GetConsoleServiceDetail(id string) {
	apiClient := client.NewAPIClient()
	detail, err := apiClient.GetConsoleServiceDetail(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching console service detail: %v\n", err)
		os.Exit(1)
	}
	printer.PrintConsoleServiceDetail(os.Stdout, detail)
}

func PublishService(req types.PublishServiceRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.PublishService(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error publishing service: %v\n", err)
		os.Exit(1)
	}
	printer.PrintPublishServiceSuccess(os.Stdout, resp)
}

func GetMyServices() {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.GetMyServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching my services: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMyServicesList(os.Stdout, resp)
}

func UpdateService(id int, req types.UpdateServiceRequest) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.UpdateService(id, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating service: %v\n", err)
		os.Exit(1)
	}
	printer.PrintUpdateServiceSuccess(os.Stdout, resp)
}

func ActivateService(id int) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.ActivateService(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error activating service: %v\n", err)
		os.Exit(1)
	}
	printer.PrintActivateServiceSuccess(os.Stdout, resp)
}

func DeactivateService(id int) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.DeactivateService(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deactivating service: %v\n", err)
		os.Exit(1)
	}
	printer.PrintDeactivateServiceSuccess(os.Stdout, resp)
}

type GetMarketServicesParams struct {
	Page     int
	Limit    int
	Search   string
	Category string
	Status   string
}

func GetMarketServiceList(params GetMarketServicesParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.MarketServicesQueryParams{
		Page:     params.Page,
		Limit:    params.Limit,
		Search:   params.Search,
		Category: params.Category,
		Status:   params.Status,
	}

	list, err := apiClient.GetMarketServices(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market services: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketServiceList(os.Stdout, list)
}

type GetMarketServicesParamsV2 struct {
	Page     int
	Limit    int
	Keyword  string
	Category string
	MinPrice float64
	MaxPrice float64
	SortBy   string
}

func GetMarketServiceListV2(params GetMarketServicesParamsV2) {
	apiClient := client.NewAPIClient()

	queryParams := client.MarketServicesQueryParamsV2{
		Page:     params.Page,
		Limit:    params.Limit,
		Keyword:  params.Keyword,
		Category: params.Category,
		MinPrice: params.MinPrice,
		MaxPrice: params.MaxPrice,
		SortBy:   params.SortBy,
	}

	list, err := apiClient.GetMarketServicesV2(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market services: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketServiceList(os.Stdout, list)
}

func GetMarketServiceDetail(id string) {
	apiClient := client.NewAPIClient()
	svc, err := apiClient.GetMarketServiceDetail(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market service detail: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketServiceDetail(os.Stdout, svc)
}

type GetServiceInvocationsParams struct {
	Role   string
	Status string
	Page   int
	Limit  int
}

func GetServiceInvocationList(params GetServiceInvocationsParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.ServiceInvocationsQueryParams{
		Role:   params.Role,
		Status: params.Status,
		Page:   params.Page,
		Limit:  params.Limit,
	}

	list, err := apiClient.GetServiceInvocations(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching service invocations: %v\n", err)
		os.Exit(1)
	}
	printer.PrintServiceInvocationList(os.Stdout, list)
}

func GetServiceInvocationDetail(id int) {
	apiClient := client.NewAPIClient()
	detail, err := apiClient.GetServiceInvocationDetail(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching service invocation detail: %v\n", err)
		os.Exit(1)
	}
	printer.PrintServiceInvocationDetail(os.Stdout, detail)
}

func InvokeService(req types.InvokeServiceRequest) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.InvokeService(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error invoking service: %v\n", err)
		os.Exit(1)
	}
	printer.PrintInvokeServiceSuccess(os.Stdout, resp)
}

func SubmitServiceResult(id int, req types.SubmitServiceResultRequest) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.SubmitServiceResult(id, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error submitting service result: %v\n", err)
		os.Exit(1)
	}
	printer.PrintSubmitServiceResultSuccess(os.Stdout, resp)
}

func ReviewServiceInvocation(id int, req types.ReviewServiceInvocationRequest) {
	apiClient := client.NewAPIClient()
	resp, err := apiClient.ReviewServiceInvocation(id, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reviewing service invocation: %v\n", err)
		os.Exit(1)
	}
	printer.PrintReviewServiceInvocationSuccess(os.Stdout, resp)
}
