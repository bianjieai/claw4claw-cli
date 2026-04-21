package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	restyClient *resty.Client
}

func NewAPIClient() *APIClient {
	client := resty.New()
	client.SetBaseURL(config.GlobalConfig.APIEndpoint)

	if config.GlobalConfig.APIToken != "" {
		client.SetHeader("X-API-Key", config.GlobalConfig.APIToken)
	}

	return &APIClient{
		restyClient: client,
	}
}

// GetMarketAgents queries public agents from the market
func (c *APIClient) GetMarketAgents() (*types.MarketAgentList, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get("/openapi/v1/market/agents")

	if err != nil {
		return nil, fmt.Errorf("failed to get market agents: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketAgentList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMarketAgentDetail queries a specific agent from the market
func (c *APIClient) GetMarketAgentDetail(id string) (*types.MarketAgent, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/market/agents/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get market agent detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketAgent
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) handleError(resp *resty.Response, env *Envelope) error {
	if env.Error != nil {
		return &APIError{StatusCode: resp.StatusCode(), Body: env.Error}
	}
	return fmt.Errorf("api request failed with status: %d", resp.StatusCode())
}

func (c *APIClient) decode(data any, target any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := json.Unmarshal(b, target); err != nil {
		return fmt.Errorf("failed to unmarshal data to target: %w", err)
	}
	return nil
}

// GetAgentStatus queries the agent status from the API.
func (c *APIClient) GetAgentStatus() (*types.AgentStatus, error) {
	var env Envelope

	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get("/api/v1/agent/status")

	if err != nil {
		return nil, fmt.Errorf("failed to call get agent status: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var status types.AgentStatus
	if err := c.decode(env.Data, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

type MarketTasksQueryParams struct {
	Page     int
	Limit    int
	Search   string
	Category string
	Status   string
}

func (c *APIClient) GetMarketTasks(params MarketTasksQueryParams) (*types.MarketTaskList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", strconv.Itoa(params.Limit))
	}
	if params.Search != "" {
		req.SetQueryParam("search", params.Search)
	}
	if params.Category != "" {
		req.SetQueryParam("category", params.Category)
	}
	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}

	resp, err := req.Get("/openapi/v1/market/tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to get market tasks: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketTaskList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMarketTaskDetail(id string) (*types.MarketTask, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/market/tasks/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get market task detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketTask
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ConsoleTasksQueryParams struct {
	Role   string
	Search string
	Status string
	SortBy string
}

func (c *APIClient) GetConsoleTasks(params ConsoleTasksQueryParams) (*types.ConsoleTaskList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Role != "" {
		req.SetQueryParam("role", params.Role)
	}
	if params.Search != "" {
		req.SetQueryParam("search", params.Search)
	}
	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}
	if params.SortBy != "" {
		req.SetQueryParam("sortBy", params.SortBy)
	}

	resp, err := req.Get("/openapi/v1/agent/me/tasks")

	if err != nil {
		return nil, fmt.Errorf("failed to get console tasks: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ConsoleTaskList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) RegisterAgent(req types.RegisterAgentReq) (*types.AgentInfo, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/register")

	if err != nil {
		return nil, fmt.Errorf("failed to register agent: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AgentInfo
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMyInfo() (*types.AgentInfo, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get("/openapi/v1/agent/me")

	if err != nil {
		return nil, fmt.Errorf("failed to get agent info: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.AgentInfo
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) UpdateMyInfo(req types.UpdateAgentReq) error {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Patch("/openapi/v1/agent/me")

	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	if resp.IsError() || !env.Success {
		return c.handleError(resp, &env)
	}

	return nil
}

func (c *APIClient) SetMyStatus(status string) error {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(types.SetStatusReq{Status: status}).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/status")

	if err != nil {
		return fmt.Errorf("failed to set agent status: %w", err)
	}

	if resp.IsError() || !env.Success {
		return c.handleError(resp, &env)
	}

	return nil
}

func (c *APIClient) UnpublishAgent() (*types.UnpublishAgentResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/unpublish")

	if err != nil {
		return nil, fmt.Errorf("failed to unpublish agent: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.UnpublishAgentResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) PublishAgent(req types.PublishAgentReq) (*types.PublishAgentResp, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/publish")

	if err != nil {
		return nil, fmt.Errorf("failed to publish agent: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.PublishAgentResp
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMyBudget() (*types.BudgetInfo, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get("/openapi/v1/agent/me/budget")

	if err != nil {
		return nil, fmt.Errorf("failed to get agent budget: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.BudgetInfo
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
