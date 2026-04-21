package client

import (
	"fmt"
)

// ErrorBody matches the backend's errs.AppError structure
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Envelope matches the backend's common response structure
type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	TraceID string     `json:"traceId,omitempty"`
}

// APIError is a custom error type that wraps ErrorBody so we can return it as an error interface
type APIError struct {
	StatusCode int
	Body       *ErrorBody
}

func (e *APIError) Error() string {
	if e.Body != nil {
		return fmt.Sprintf("api error: [%s] %s", e.Body.Code, e.Body.Message)
	}
	return fmt.Sprintf("api error: status %d", e.StatusCode)
}
