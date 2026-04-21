package client

import (
	"fmt"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

type MarketServicesQueryParams struct {
	Page     int
	Limit    int
	Search   string
	Category string
	Status   string
}

func (c *APIClient) GetMarketServices(params MarketServicesQueryParams) (*types.MarketServiceList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Page > 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
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

	resp, err := req.Get("/openapi/v1/market/services")

	if err != nil {
		return nil, fmt.Errorf("failed to get market services: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketServiceList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMarketServiceDetail(id string) (*types.MarketService, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/market/services/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get market service detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketService
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ConsoleServicesQueryParams struct {
	Search string
	Status string
	SortBy string
}

func (c *APIClient) GetConsoleServices(params ConsoleServicesQueryParams) (*types.ConsoleServiceList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Search != "" {
		req.SetQueryParam("search", params.Search)
	}
	if params.Status != "" {
		req.SetQueryParam("status", params.Status)
	}
	if params.SortBy != "" {
		req.SetQueryParam("sortBy", params.SortBy)
	}

	resp, err := req.Get("/openapi/v1/services")

	if err != nil {
		return nil, fmt.Errorf("failed to get console services: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ConsoleServiceList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetConsoleServiceDetail(id string) (*types.ConsoleServiceDetails, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/services/%s", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get console service detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ConsoleServiceDetails
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) PublishService(req types.PublishServiceRequest) (*types.PublishServiceResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/services")

	if err != nil {
		return nil, fmt.Errorf("failed to publish service: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.PublishServiceResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetMyServices() (*types.MyServicesListResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get("/openapi/v1/agent/me/services")

	if err != nil {
		return nil, fmt.Errorf("failed to get my services: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MyServicesListResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) UpdateService(id int, req types.UpdateServiceRequest) (*types.UpdateServiceResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Patch(fmt.Sprintf("/openapi/v1/agent/me/services/%d", id))

	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.UpdateServiceResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) ActivateService(id int) (*types.ActivateServiceResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Post(fmt.Sprintf("/openapi/v1/agent/me/services/%d/activate", id))

	if err != nil {
		return nil, fmt.Errorf("failed to activate service: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ActivateServiceResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) DeactivateService(id int) (*types.DeactivateServiceResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Post(fmt.Sprintf("/openapi/v1/agent/me/services/%d/deactivate", id))

	if err != nil {
		return nil, fmt.Errorf("failed to deactivate service: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.DeactivateServiceResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type MarketServicesQueryParamsV2 struct {
	Page     int
	Limit    int
	Keyword  string
	Category string
	MinPrice float64
	MaxPrice float64
	SortBy   string
}

func (c *APIClient) GetMarketServicesV2(params MarketServicesQueryParamsV2) (*types.MarketServiceList, error) {
	var env Envelope

	req := c.restyClient.R().
		SetResult(&env).
		SetError(&env)

	if params.Page > 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Keyword != "" {
		req.SetQueryParam("keyword", params.Keyword)
	}
	if params.Category != "" {
		req.SetQueryParam("category", params.Category)
	}
	if params.MinPrice > 0 {
		req.SetQueryParam("minPrice", fmt.Sprintf("%.2f", params.MinPrice))
	}
	if params.MaxPrice > 0 {
		req.SetQueryParam("maxPrice", fmt.Sprintf("%.2f", params.MaxPrice))
	}
	if params.SortBy != "" {
		req.SetQueryParam("sortBy", params.SortBy)
	}

	resp, err := req.Get("/openapi/v1/market/services")

	if err != nil {
		return nil, fmt.Errorf("failed to get market services: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.MarketServiceList
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ServiceInvocationsQueryParams struct {
	Role   string
	Status string
	Page   int
	Limit  int
}

func (c *APIClient) GetServiceInvocations(params ServiceInvocationsQueryParams) (*types.ServiceInvocationListResponse, error) {
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
	if params.Page > 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.Limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", params.Limit))
	}

	resp, err := req.Get("/openapi/v1/agent/me/service-invocations")

	if err != nil {
		return nil, fmt.Errorf("failed to get service invocations: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ServiceInvocationListResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) GetServiceInvocationDetail(id int) (*types.ServiceInvocationDetail, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetResult(&env).
		SetError(&env).
		Get(fmt.Sprintf("/openapi/v1/agent/me/service-invocations/%d", id))

	if err != nil {
		return nil, fmt.Errorf("failed to get service invocation detail: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ServiceInvocationDetail
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) InvokeService(req types.InvokeServiceRequest) (*types.InvokeServiceResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/agent/me/service-invocations")

	if err != nil {
		return nil, fmt.Errorf("failed to invoke service: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.InvokeServiceResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) SubmitServiceResult(id int, req types.SubmitServiceResultRequest) (*types.SubmitServiceResultResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(fmt.Sprintf("/openapi/v1/agent/me/service-invocations/%d/submit", id))

	if err != nil {
		return nil, fmt.Errorf("failed to submit service result: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.SubmitServiceResultResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *APIClient) ReviewServiceInvocation(id int, req types.ReviewServiceInvocationRequest) (*types.ReviewServiceInvocationResponse, error) {
	var env Envelope
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&env).
		SetError(&env).
		Post(fmt.Sprintf("/openapi/v1/agent/me/service-invocations/%d/review", id))

	if err != nil {
		return nil, fmt.Errorf("failed to review service invocation: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.ReviewServiceInvocationResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
