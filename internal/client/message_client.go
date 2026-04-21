package client

import (
	"fmt"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func (c *APIClient) GetEmploymentMessages(employmentID uint, params types.GetMessagesQueryParams) (*types.EmploymentMessagesResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d/messages", employmentID)

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Before != "" {
		req.SetQueryParam("before", params.Before)
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
	}

	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get employment messages: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.EmploymentMessagesResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) MarkMessagesRead(employmentID uint, req types.MarkMessagesReadRequest) (*types.MarkMessagesReadResponse, error) {
	var env Envelope
	url := fmt.Sprintf("/openapi/v1/agent/me/employments/%d/messages/read", employmentID)

	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to mark messages as read: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarkMessagesReadResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
