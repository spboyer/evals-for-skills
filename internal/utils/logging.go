package utils

import (
	"context"
	"log/slog"

	copilot "github.com/github/copilot-sdk/go"
)

func SessionToSlog(event copilot.SessionEvent) {
	if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		return
	}

	attrs := []any{
		"type", event.Type,
	}

	attrs = addIf(attrs, "content", event.Data.Content)
	attrs = addIf(attrs, "deltaContent", event.Data.DeltaContent)
	attrs = addIf(attrs, "toolName", event.Data.ToolName)
	attrs = addIf(attrs, "toolResult", event.Data.Result)
	attrs = addIf(attrs, "toolCallID", event.Data.ToolCallID)
	attrs = addIf(attrs, "reasoningText", event.Data.ReasoningText)

	slog.Debug("Event received", attrs...)
}

func addIf[T any](attrs []any, name string, v *T) []any {
	if v != nil {
		attrs = append(attrs, name)
		attrs = append(attrs, *v)
	}

	return attrs
}
