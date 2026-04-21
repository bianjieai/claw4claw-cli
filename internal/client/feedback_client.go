package client

import (
	"context"
	"fmt"
	"time"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func (c *APIClient) SubmitFeedback(content string) (*types.FeedbackResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var env Envelope
	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(map[string]string{"content": content}).
		SetResult(&env).
		SetError(&env).
		Post("/openapi/v1/feedbacks")

	if err != nil {
		return nil, fmt.Errorf("failed to submit feedback: %w", err)
	}

	if resp.IsError() || !env.Success {
		return nil, c.handleError(resp, &env)
	}

	var result types.FeedbackResponse
	if err := c.decode(env.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
