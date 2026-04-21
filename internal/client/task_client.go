package client

import (
	"fmt"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func (c *APIClient) PublishTask(req types.PublishTaskRequest) (*types.PublishTaskResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to publish task: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.PublishTaskResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type MyTasksQueryParams struct {
	Role   string
	Status string
	Search string
	SortBy string
}

func (c *APIClient) GetMyTasks(params MyTasksQueryParams) (*types.MyTaskList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Role != "" {
		req.SetQueryParam("role", params.Role)
	}
	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}
	if params.Search != "" {
		req.SetQueryParam("search", params.Search)
	}
	if params.SortBy != "" {
		req.SetQueryParam("sortBy", params.SortBy)
	}

	resp, err := req.Get("/openapi/v1/agent/me/tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to get my tasks: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MyTaskList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetTaskDetail(taskID string) (*types.TaskDetail, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/agent/me/tasks/%s", taskID))

	if err != nil {
		return nil, fmt.Errorf("failed to get task detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.TaskDetail
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetTaskReview(taskID string) (*types.TaskReview, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/agent/me/tasks/%s", taskID))

	if err != nil {
		return nil, fmt.Errorf("failed to get task review: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.TaskReview
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) ApplyForTask(req types.ApplyTaskRequest) (*types.ApplyTaskResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/task-applications")

	if err != nil {
		return nil, fmt.Errorf("failed to apply for task: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ApplyTaskResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type MyTaskApplicationsQueryParams struct {
	Status string
	Page   int
	Limit  int
}

func (c *APIClient) GetMyTaskApplications(params MyTaskApplicationsQueryParams) (*types.MyTaskApplicationList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}
	if params.Page > 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
	}

	resp, err := req.Get("/openapi/v1/agent/me/task-applications")

	if err != nil {
		return nil, fmt.Errorf("failed to get my task applications: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MyTaskApplicationList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) WithdrawTaskApplication(applicationID string) (*types.WithdrawApplicationResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Delete(fmt.Sprintf("/openapi/v1/agent/me/task-applications/%s", applicationID))

	if err != nil {
		return nil, fmt.Errorf("failed to withdraw task application: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.WithdrawApplicationResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) SubmitTaskResult(applicationID string, req types.SubmitTaskRequest) (*types.SubmitTaskResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/task-applications/%s/submit", applicationID)
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to submit task result: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.SubmitTaskResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type AcceptedTasksQueryParams struct {
	Status string
}

func (c *APIClient) GetAcceptedTasks(params AcceptedTasksQueryParams) (*types.AcceptedTaskList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}

	resp, err := req.Get("/openapi/v1/agent/me/task-applications/accepted")

	if err != nil {
		return nil, fmt.Errorf("failed to get accepted tasks: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AcceptedTaskList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type TaskApplicationsQueryParams struct {
	Status string
}

func (c *APIClient) GetTaskApplications(taskID string, params TaskApplicationsQueryParams) (*types.TaskApplicationList, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/tasks/%s/applications", taskID)

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Status != "" {
		req = req.SetQueryParam("status", params.Status)
	}

	resp, err := req.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get task applications: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.TaskApplicationList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) AcceptApplicant(taskID string, applicationID string, req types.AcceptApplicantRequest) (*types.AcceptApplicantResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/tasks/%s/applications/%s/accept", taskID, applicationID)
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to accept applicant: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AcceptApplicantResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) AcceptTask(taskID string, req types.AcceptTaskRequest) (*types.AcceptTaskResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/tasks/%s/accept", taskID)
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to accept task: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AcceptTaskResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) CancelTask(taskID string) (*types.CancelTaskResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/tasks/%s/cancel", taskID)
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to cancel task: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.CancelTaskResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
