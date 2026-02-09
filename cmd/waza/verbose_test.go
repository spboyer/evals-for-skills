package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spboyer/waza/internal/orchestration"
)

// captureOutput captures stdout during function execution.
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestVerboseProgressListener_BenchmarkStart(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventBenchmarkStart,
			TotalTests: 5,
		})
	})
	if !strings.Contains(out, "5 test(s)") {
		t.Errorf("expected test count in output, got: %s", out)
	}
}

func TestVerboseProgressListener_TestStart(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventTestStart,
			TestName:   "my-test",
			TestNum:    2,
			TotalTests: 3,
		})
	})
	if !strings.Contains(out, "[2/3]") {
		t.Errorf("expected [2/3] in output, got: %s", out)
	}
	if !strings.Contains(out, "my-test") {
		t.Errorf("expected test name in output, got: %s", out)
	}
}

func TestVerboseProgressListener_AgentPrompt(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentPrompt,
			TestName:  "test-1",
			Details:   map[string]any{"message": "Explain this code"},
		})
	})
	if !strings.Contains(out, "Agent Prompt") {
		t.Errorf("expected 'Agent Prompt' header, got: %s", out)
	}
	if !strings.Contains(out, "Explain this code") {
		t.Errorf("expected prompt message in output, got: %s", out)
	}
}

func TestVerboseProgressListener_AgentResponse(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			TestName:  "test-1",
			Details: map[string]any{
				"output":     "This code does X\nand also Y",
				"tool_calls": 3,
			},
		})
	})
	if !strings.Contains(out, "Agent Response") {
		t.Errorf("expected 'Agent Response' header, got: %s", out)
	}
	if !strings.Contains(out, "This code does X") {
		t.Errorf("expected response output, got: %s", out)
	}
	if !strings.Contains(out, "Tool calls: 3") {
		t.Errorf("expected tool call count, got: %s", out)
	}
}

func TestVerboseProgressListener_AgentResponse_NoToolCalls(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType: orchestration.EventAgentResponse,
			TestName:  "test-1",
			Details: map[string]any{
				"output":     "Simple response",
				"tool_calls": 0,
			},
		})
	})
	if strings.Contains(out, "Tool calls") {
		t.Errorf("should not show tool calls when zero, got: %s", out)
	}
}

func TestVerboseProgressListener_GraderResult_Passed(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventGraderResult,
			TestName:   "test-1",
			DurationMs: 42,
			Details: map[string]any{
				"grader":   "regex-check",
				"type":     "regex",
				"passed":   true,
				"score":    1.0,
				"feedback": "All patterns matched",
			},
		})
	})
	if !strings.Contains(out, "✓") {
		t.Errorf("expected pass icon, got: %s", out)
	}
	if !strings.Contains(out, "regex-check") {
		t.Errorf("expected grader name, got: %s", out)
	}
	if !strings.Contains(out, "score=1.00") {
		t.Errorf("expected score, got: %s", out)
	}
	if !strings.Contains(out, "42ms") {
		t.Errorf("expected duration, got: %s", out)
	}
	// Feedback should NOT appear for passed graders
	if strings.Contains(out, "Feedback:") {
		t.Errorf("should not show feedback for passing grader, got: %s", out)
	}
}

func TestVerboseProgressListener_GraderResult_Failed(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventGraderResult,
			TestName:   "test-1",
			DurationMs: 15,
			Details: map[string]any{
				"grader":   "code-check",
				"type":     "code",
				"passed":   false,
				"score":    0.5,
				"feedback": "Missing function signature",
			},
		})
	})
	if !strings.Contains(out, "✗") {
		t.Errorf("expected fail icon, got: %s", out)
	}
	if !strings.Contains(out, "score=0.50") {
		t.Errorf("expected score, got: %s", out)
	}
	if !strings.Contains(out, "Feedback: Missing function signature") {
		t.Errorf("expected feedback for failed grader, got: %s", out)
	}
}

func TestVerboseProgressListener_RunComplete(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventRunComplete,
			Status:     "passed",
			DurationMs: 1500,
		})
	})
	duration := time.Duration(1500) * time.Millisecond
	expected := fmt.Sprintf("Run result: passed (%v)", duration)
	if !strings.Contains(out, expected) {
		t.Errorf("expected %q in output, got: %s", expected, out)
	}
}

func TestVerboseProgressListener_BenchmarkComplete(t *testing.T) {
	out := captureOutput(func() {
		verboseProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventBenchmarkComplete,
			DurationMs: 5000,
		})
	})
	if !strings.Contains(out, "Benchmark completed") {
		t.Errorf("expected benchmark complete message, got: %s", out)
	}
}

func TestSimpleProgressListener_PassedTest(t *testing.T) {
	out := captureOutput(func() {
		simpleProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventTestComplete,
			TestName:   "simple-test",
			TestNum:    1,
			TotalTests: 2,
			Status:     "passed",
		})
	})
	if !strings.Contains(out, "✓") {
		t.Errorf("expected pass icon, got: %s", out)
	}
	if !strings.Contains(out, "simple-test") {
		t.Errorf("expected test name, got: %s", out)
	}
}

func TestSimpleProgressListener_FailedTest(t *testing.T) {
	out := captureOutput(func() {
		simpleProgressListener(orchestration.ProgressEvent{
			EventType:  orchestration.EventTestComplete,
			TestName:   "fail-test",
			TestNum:    1,
			TotalTests: 1,
			Status:     "failed",
		})
	})
	if !strings.Contains(out, "✗") {
		t.Errorf("expected fail icon, got: %s", out)
	}
}
