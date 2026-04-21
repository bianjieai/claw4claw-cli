package service

import (
	"fmt"
	"os"

	"github.com/bianjieai/claw4claw-cli/internal/client"
	"github.com/bianjieai/claw4claw-cli/internal/printer"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

type GetMarketTasksParams struct {
	Page     int
	Limit    int
	Search   string
	Category string
	Status   string
}

func GetMarketTaskList(params GetMarketTasksParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.MarketTasksQueryParams{
		Page:     params.Page,
		Limit:    params.Limit,
		Search:   params.Search,
		Category: params.Category,
		Status:   params.Status,
	}

	list, err := apiClient.GetMarketTasks(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market tasks: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketTaskList(os.Stdout, list)
}

func GetMarketTaskDetail(id string) {
	apiClient := client.NewAPIClient()
	task, err := apiClient.GetMarketTaskDetail(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching market task detail: %v\n", err)
		os.Exit(1)
	}
	printer.PrintMarketTaskDetail(os.Stdout, task)
}

type GetConsoleTasksParams struct {
	Role   string
	Search string
	Status string
	SortBy string
}

func GetConsoleTaskList(params GetConsoleTasksParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.ConsoleTasksQueryParams{
		Role:   params.Role,
		Search: params.Search,
		Status: params.Status,
		SortBy: params.SortBy,
	}

	list, err := apiClient.GetConsoleTasks(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching console tasks: %v\n", err)
		os.Exit(1)
	}
	printer.PrintConsoleTaskList(os.Stdout, list, params.Role)
}

func PublishTask(req types.PublishTaskRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.PublishTask(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error publishing task: %v\n", err)
		os.Exit(1)
	}
	printer.PrintPublishTaskSuccess(os.Stdout, resp)
}

func ApplyTask(taskID string, req types.ApplyTaskRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.ApplyForTask(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error applying task: %v\n", err)
		os.Exit(1)
	}
	printer.PrintApplyTaskSuccess(os.Stdout, resp)
}

func SubmitTask(applicationID string, req types.SubmitTaskRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.SubmitTaskResult(applicationID, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error submitting task: %v\n", err)
		os.Exit(1)
	}
	printer.PrintSubmitTaskSuccess(os.Stdout, resp)
}

func AcceptTask(taskID string, req types.AcceptTaskRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.AcceptTask(taskID, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accepting task: %v\n", err)
		os.Exit(1)
	}
	printer.PrintAcceptTaskSuccess(os.Stdout, resp)
}

type GetTaskApplicationsParams struct {
	Status string
}

func GetTaskApplications(taskID string, params GetTaskApplicationsParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.TaskApplicationsQueryParams{
		Status: params.Status,
	}

	list, err := apiClient.GetTaskApplications(taskID, queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching task applications: %v\n", err)
		os.Exit(1)
	}
	printer.PrintTaskApplicationList(os.Stdout, list, taskID)
}

func AcceptApplicant(taskID string, applicationID string, req types.AcceptApplicantRequest) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.AcceptApplicant(taskID, applicationID, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accepting applicant: %v\n", err)
		os.Exit(1)
	}
	printer.PrintAcceptApplicantSuccess(os.Stdout, resp)
}

func CancelTask(taskID string) {
	apiClient := client.NewAPIClient()

	resp, err := apiClient.CancelTask(taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cancelling task: %v\n", err)
		os.Exit(1)
	}
	printer.PrintCancelTaskSuccess(os.Stdout, resp)
}

type GetAcceptedTasksParams struct {
	Status string
}

func GetAcceptedTasks(params GetAcceptedTasksParams) {
	apiClient := client.NewAPIClient()

	queryParams := client.AcceptedTasksQueryParams{
		Status: params.Status,
	}

	list, err := apiClient.GetAcceptedTasks(queryParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching accepted tasks: %v\n", err)
		os.Exit(1)
	}
	printer.PrintAcceptedTaskList(os.Stdout, list)
}
