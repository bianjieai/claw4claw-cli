package printer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func TestFormatNotificationEvent(t *testing.T) {
	tests := []struct {
		name     string
		notif    types.WebSocketNotificationMessage
		contains string
	}{
		{
			name: "task application",
			notif: types.WebSocketNotificationMessage{
				Domain: "task",
				Event:  "task_application",
				Data:   json.RawMessage(`{"taskTitle":"Task A","applicantAgentName":"agent-x","bounty":"10"}`),
			},
			contains: "applied for your task",
		},
		{
			name: "service timeout",
			notif: types.WebSocketNotificationMessage{
				Domain: "service_invocation",
				Event:  "invocation_timeout",
				Data:   json.RawMessage(`{"serviceTitle":"svc-a"}`),
			},
			contains: "timed out",
		},
		{
			name: "employment completed",
			notif: types.WebSocketNotificationMessage{
				Domain: "employment",
				Event:  "employment_completed",
				Data:   json.RawMessage(`{"employerAgentName":"boss"}`),
			},
			contains: "has completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatNotificationEvent(tt.notif)
			if !strings.Contains(got, tt.contains) {
				t.Fatalf("expected %q to contain %q", got, tt.contains)
			}
		})
	}
}
