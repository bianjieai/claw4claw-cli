package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bianjieai/claw4claw-cli/internal/config"
	"github.com/bianjieai/claw4claw-cli/internal/types"
)

func TestPrintAgentStatus_JSON(t *testing.T) {
	config.GlobalConfig.OutputFormat = "json"

	data := &types.AgentStatus{
		DID:     "did:c4c:test",
		Name:    "TestAgent",
		Status:  "online",
		Balance: "100 C4C",
		Reputation: types.Reputation{
			Level: "Pro",
			Score: 90.0,
		},
	}

	var buf bytes.Buffer
	PrintAgentStatus(&buf, data)

	output := buf.String()
	if !strings.Contains(output, `"did": "did:c4c:test"`) {
		t.Errorf("Expected output to contain DID, got: %s", output)
	}
	if !strings.Contains(output, `"name": "TestAgent"`) {
		t.Errorf("Expected output to contain Name, got: %s", output)
	}
}

func TestPrintAgentStatus_Text(t *testing.T) {
	config.GlobalConfig.OutputFormat = "text"

	data := &types.AgentStatus{
		DID:     "did:c4c:test",
		Name:    "TestAgent",
		Status:  "online",
		Balance: "100 C4C",
		Reputation: types.Reputation{
			Level: "Pro",
			Score: 90.0,
		},
	}

	var buf bytes.Buffer
	PrintAgentStatus(&buf, data)

	output := buf.String()
	if !strings.Contains(output, "Agent Status") {
		t.Errorf("Expected output to contain header, got: %s", output)
	}
	if !strings.Contains(output, "did:c4c:test") {
		t.Errorf("Expected output to contain DID, got: %s", output)
	}
}
