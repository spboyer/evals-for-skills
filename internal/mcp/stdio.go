package mcp

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/spboyer/waza/internal/jsonrpc"
)

// ServeStdio runs the MCP server on the given reader/writer (typically stdin/stdout).
// It reads newline-delimited JSON-RPC requests and writes responses.
func ServeStdio(ctx context.Context, r io.Reader, w io.Writer, logger *slog.Logger) {
	srv := NewServer(logger)
	transport := jsonrpc.NewTransport(r, w)

	for {
		req, rawJSON, err := transport.ReadRequest()
		if err != nil {
			if err == io.EOF {
				return
			}
			logger.Debug("mcp read error", "error", err)
			resp := &jsonrpc.Response{
				JSONRPC: "2.0",
				Error:   jsonrpc.ErrParseError(err.Error()),
				ID:      json.RawMessage("null"),
			}
			if writeErr := transport.WriteResponse(resp); writeErr != nil {
				logger.Debug("mcp write error", "error", writeErr)
			}
			return
		}

		// Detect notifications (no "id" field).
		isNotification := !hasIDField(rawJSON)

		resp := srv.HandleRequest(ctx, req)

		// MCP initialize sends notifications/initialized as follow-up — skip response.
		if resp == nil || isNotification {
			continue
		}

		if writeErr := transport.WriteResponse(resp); writeErr != nil {
			logger.Debug("mcp write error", "error", writeErr)
			return
		}
	}
}

// hasIDField checks whether the raw JSON contains an "id" key at the top level.
func hasIDField(raw []byte) bool {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return false
	}
	_, exists := obj["id"]
	return exists
}
