package utils

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	copilot "github.com/github/copilot-sdk/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionToSlogDebugDisabled(t *testing.T) {
	old := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(old)
	})

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	SessionToSlog(copilot.SessionEvent{Type: copilot.SessionEventType("message")})
	assert.Equal(t, 0, buf.Len())
}

func TestSessionToSlogDebugEnabled(t *testing.T) {
	old := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(old)
	})

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	content := "hello"
	deltaContent := " world"
	toolName := "bash"
	toolCallID := "call-1"
	reasoningText := "reasoning"

	SessionToSlog(copilot.SessionEvent{
		Type: copilot.SessionEventType("message"),
		Data: copilot.Data{
			Content:       &content,
			DeltaContent:  &deltaContent,
			ToolName:      &toolName,
			ToolCallID:    &toolCallID,
			ReasoningText: &reasoningText,
		},
	})

	var logEntry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))
	assert.Equal(t, "Event received", logEntry["msg"])
	assert.Equal(t, "message", logEntry["type"])
	assert.Equal(t, content, logEntry["content"])
	assert.Equal(t, deltaContent, logEntry["deltaContent"])
	assert.Equal(t, toolName, logEntry["toolName"])
	assert.Equal(t, toolCallID, logEntry["toolCallID"])
	assert.Equal(t, reasoningText, logEntry["reasoningText"])
}

func TestAddIf(t *testing.T) {
	attrs := []any{"existing", "value"}

	result := addIf(attrs, "missing", (*int)(nil))
	assert.Equal(t, attrs, result)

	v := 7
	result = addIf(attrs, "number", &v)
	assert.Equal(t, []any{"existing", "value", "number", 7}, result)
}
