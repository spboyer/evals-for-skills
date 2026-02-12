package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spboyer/waza/internal/orchestration"
	"github.com/stretchr/testify/assert"
)

// captureStdout redirects os.Stdout and returns captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return buf.String()
}

func TestVerbose_AgentPrompt(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentPrompt,
			Details:   map[string]any{"message": "Explain this code"},
		})
	})
	assert.Contains(t, out, "[PROMPT]")
	assert.Contains(t, out, "Explain this code")
}

func TestVerbose_AgentResponse(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"output":     "This is the response",
				"tool_calls": 3,
			},
		})
	})
	assert.Contains(t, out, "[RESPONSE]")
	assert.Contains(t, out, "This is the response")
	assert.Contains(t, out, "[TOOLS] 3 tool call(s)")
}

func TestVerbose_AgentResponse_NoToolCalls(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"output":     "short answer",
				"tool_calls": 0,
			},
		})
	})
	assert.Contains(t, out, "[RESPONSE]")
	assert.NotContains(t, out, "[TOOLS]")
}

func TestVerbose_GraderResult_Passed(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventGraderResult,
			DurationMs: 42,
			Details: map[string]any{
				"grader":   "contains_check",
				"passed":   true,
				"score":    1.0,
				"feedback": "",
			},
		})
	})
	assert.Contains(t, out, "✓")
	assert.Contains(t, out, "contains_check")
	assert.Contains(t, out, "score=1.00")
}

func TestVerbose_GraderResult_Failed(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventGraderResult,
			DurationMs: 15,
			Details: map[string]any{
				"grader":   "regex_match",
				"passed":   false,
				"score":    0.0,
				"feedback": "pattern not found",
			},
		})
	})
	assert.Contains(t, out, "✗")
	assert.Contains(t, out, "regex_match")
	assert.Contains(t, out, "pattern not found")
}

func TestVerbose_AgentResponse_ErrorMsg(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"error":      "copilot session crashed",
				"output":     "",
				"tool_calls": 0,
			},
		})
	})
	assert.Contains(t, out, "[ERROR]")
	assert.Contains(t, out, "copilot session crashed")
	assert.NotContains(t, out, "[RESPONSE]")
}

func TestVerbose_AgentResponse_EmptyErrorMsg(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"error":      "",
				"output":     "ok",
				"tool_calls": 0,
			},
		})
	})
	assert.NotContains(t, out, "[ERROR]")
}

func TestVerbose_AgentResponse_NoErrorMsgKey(t *testing.T) {
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"output":     "ok",
				"tool_calls": 0,
			},
		})
	})
	assert.NotContains(t, out, "[ERROR]")
	assert.Contains(t, out, "[RESPONSE]")
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is a ..."},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("len=%d", tt.maxLen), func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVerbose_ResponseTruncation(t *testing.T) {
	longOutput := strings.Repeat("x", 300)
	out := captureStdout(t, func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			Details: map[string]any{
				"output":     longOutput,
				"tool_calls": 0,
			},
		})
	})
	assert.Contains(t, out, "...")
	// Truncated to 200 + "..." = should not contain full 300 chars
	assert.Less(t, len(out), 300)
}
