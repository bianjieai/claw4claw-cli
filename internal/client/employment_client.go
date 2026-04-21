package client

import (
	"context"
	"fmt"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func (c *APIClient) CreateEmployment(req types.CreateEmploymentRequest) (*types.CreateEmploymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/employments")

	if err != nil {
		return nil, fmt.Errorf("failed to create employment: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.CreateEmploymentResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) AcceptEmployment(employmentID uint, req types.AcceptEmploymentRequest) (*types.AcceptEmploymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d/accept", employmentID)
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to accept employment: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AcceptEmploymentResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) RejectEmployment(employmentID uint, req types.RejectEmploymentRequest) (*types.RejectEmploymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d/reject", employmentID)
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to reject employment: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.RejectEmploymentResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) TerminateEmployment(employmentID uint, req types.TerminateEmploymentRequest) (*types.TerminateEmploymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d/terminate", employmentID)
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to terminate employment: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.TerminateEmploymentResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMyEmployments(params types.MyEmploymentsQueryParams) (*types.MyEmploymentsListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope

	req := c.restyClient.R().
		SetContext(ctx).
		SetResult(&env).
		SetError(&env)

	if params.Role != "" {
		req.SetQueryParam("role", params.Role)
	}
	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}
	if params.Page > 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
	}

	resp, err := req.Get("/openapi/v1/agent/me/employments")

	if err != nil {
		return nil, fmt.Errorf("failed to get my employments: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MyEmploymentsListResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetEmploymentByID(employmentID uint) (*types.EmploymentDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d", employmentID)
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetResult(&env).
		SetError(&env).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get employment detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.EmploymentDetail
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
