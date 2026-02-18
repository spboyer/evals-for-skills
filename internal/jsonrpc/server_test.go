package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_MethodNotFound(t *testing.T) {
	registry := NewMethodRegistry()
	server := NewServer(registry, nil)

	req := `{"jsonrpc":"2.0","method":"nonexistent","id":1}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, CodeMethodNotFound, resp.Error.Code)
}

func TestServer_InvalidJSON(t *testing.T) {
	registry := NewMethodRegistry()
	server := NewServer(registry, nil)

	req := `{not valid json}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Equal(t, CodeParseError, resp.Error.Code)
}

func TestServer_InvalidVersion(t *testing.T) {
	registry := NewMethodRegistry()
	server := NewServer(registry, nil)

	req := `{"jsonrpc":"1.0","method":"test","id":1}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Equal(t, CodeInvalidRequest, resp.Error.Code)
}

func TestServer_SuccessfulMethod(t *testing.T) {
	registry := NewMethodRegistry()
	registry.Register("echo", func(_ context.Context, params json.RawMessage) (any, *Error) {
		return params, nil
	})
	server := NewServer(registry, nil)

	req := `{"jsonrpc":"2.0","method":"echo","params":{"hello":"world"},"id":42}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
}

func TestServer_MultipleRequests(t *testing.T) {
	registry := NewMethodRegistry()
	callCount := 0
	registry.Register("ping", func(_ context.Context, _ json.RawMessage) (any, *Error) {
		callCount++
		return map[string]string{"pong": "ok"}, nil
	})
	server := NewServer(registry, nil)

	reqs := `{"jsonrpc":"2.0","method":"ping","id":1}` + "\n" +
		`{"jsonrpc":"2.0","method":"ping","id":2}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(reqs), &out)

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, 2, callCount)

	for i, line := range lines {
		var resp Response
		require.NoError(t, json.Unmarshal([]byte(line), &resp), "line %d", i)
		assert.Nil(t, resp.Error)
	}
}

func TestServer_Notification_NoResponse(t *testing.T) {
	registry := NewMethodRegistry()
	called := false
	registry.Register("notify.test", func(_ context.Context, _ json.RawMessage) (any, *Error) {
		called = true
		return map[string]string{"ok": "true"}, nil
	})
	server := NewServer(registry, nil)

	// Notification: no "id" field at all
	notif := `{"jsonrpc":"2.0","method":"notify.test","params":{}}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(notif), &out)

	assert.True(t, called, "handler should be invoked for notifications")
	assert.Empty(t, out.String(), "no response should be written for notifications")
}

func TestServer_Notification_WithNullID_GetsResponse(t *testing.T) {
	registry := NewMethodRegistry()
	registry.Register("echo", func(_ context.Context, params json.RawMessage) (any, *Error) {
		return params, nil
	})
	server := NewServer(registry, nil)

	// Explicit "id": null is a request, not a notification — should get a response.
	req := `{"jsonrpc":"2.0","method":"echo","params":{"x":1},"id":null}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	assert.NotEmpty(t, out.String(), "request with explicit id:null should get a response")
	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Nil(t, resp.Error)
}

func TestServer_Notification_MethodNotFound_NoResponse(t *testing.T) {
	registry := NewMethodRegistry()
	server := NewServer(registry, nil)

	// Notification to nonexistent method: no response per spec
	notif := `{"jsonrpc":"2.0","method":"nonexistent"}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(notif), &out)

	assert.Empty(t, out.String(), "no response for notification even if method not found")
}

func TestServer_MixedRequestsAndNotifications(t *testing.T) {
	registry := NewMethodRegistry()
	callCount := 0
	registry.Register("ping", func(_ context.Context, _ json.RawMessage) (any, *Error) {
		callCount++
		return map[string]string{"pong": "ok"}, nil
	})
	server := NewServer(registry, nil)

	// Mix of request (has id) and notification (no id)
	input := `{"jsonrpc":"2.0","method":"ping","id":1}` + "\n" +
		`{"jsonrpc":"2.0","method":"ping"}` + "\n" +
		`{"jsonrpc":"2.0","method":"ping","id":2}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(input), &out)

	assert.Equal(t, 3, callCount, "all three calls should be invoked")

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	assert.Equal(t, 2, len(lines), "only 2 responses for 2 requests (notification gets none)")
}

func TestServer_ErrorFromHandler(t *testing.T) {
	registry := NewMethodRegistry()
	registry.Register("fail", func(_ context.Context, _ json.RawMessage) (any, *Error) {
		return nil, ErrInternalError("something broke")
	})
	server := NewServer(registry, nil)

	req := `{"jsonrpc":"2.0","method":"fail","id":1}` + "\n"
	var out bytes.Buffer
	server.ServeStdio(strings.NewReader(req), &out)

	var resp Response
	require.NoError(t, json.Unmarshal(out.Bytes(), &resp))
	assert.Equal(t, CodeInternalError, resp.Error.Code)
	assert.Equal(t, "Internal error", resp.Error.Message)
}

func TestTransport_WriteNotification(t *testing.T) {
	var buf bytes.Buffer
	transport := NewTransport(strings.NewReader(""), &buf)

	notif := &Notification{
		JSONRPC: "2.0",
		Method:  "eval.progress",
		Params:  map[string]any{"task": "test-1", "progress": 50},
	}

	require.NoError(t, transport.WriteNotification(notif))

	var decoded Notification
	require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
	assert.Equal(t, "2.0", decoded.JSONRPC)
	assert.Equal(t, "eval.progress", decoded.Method)
}

func TestMethodRegistry(t *testing.T) {
	reg := NewMethodRegistry()

	assert.Nil(t, reg.Lookup("test"))
	assert.Empty(t, reg.Methods())

	reg.Register("test.method", func(_ context.Context, _ json.RawMessage) (any, *Error) {
		return nil, nil
	})

	assert.NotNil(t, reg.Lookup("test.method"))
	assert.Nil(t, reg.Lookup("other"))
	assert.Equal(t, []string{"test.method"}, reg.Methods())
}

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		err  *Error
		code int
		msg  string
	}{
		{ErrParseError("bad"), CodeParseError, "Parse error"},
		{ErrInvalidRequest("bad"), CodeInvalidRequest, "Invalid request"},
		{ErrMethodNotFound("x"), CodeMethodNotFound, "Method not found"},
		{ErrInvalidParams("bad"), CodeInvalidParams, "Invalid params"},
		{ErrInternalError("bad"), CodeInternalError, "Internal error"},
		{ErrEvalNotFound("x"), CodeEvalNotFound, "Eval not found"},
		{ErrValidationFailed("bad"), CodeValidationFailed, "Validation failed"},
		{ErrRunFailed("bad"), CodeRunFailed, "Run failed"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.code, tt.err.Code)
		assert.Equal(t, tt.msg, tt.err.Message)
		assert.Equal(t, tt.msg, tt.err.Error())
	}
}
