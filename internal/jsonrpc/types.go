package jsonrpc

import "encoding/json"

// JSON-RPC 2.0 types per https://www.jsonrpc.org/specification

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      json.RawMessage `json:"id"`
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  any             `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

// Notification represents a server-initiated JSON-RPC 2.0 notification (no ID).
type Notification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// Error represents a JSON-RPC 2.0 error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// Standard JSON-RPC 2.0 error codes.
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)

// Application-specific error codes.
const (
	CodeEvalNotFound     = -32000
	CodeValidationFailed = -32001
	CodeRunFailed        = -32002
)

func ErrParseError(data any) *Error {
	return &Error{Code: CodeParseError, Message: "Parse error", Data: data}
}

func ErrInvalidRequest(data any) *Error {
	return &Error{Code: CodeInvalidRequest, Message: "Invalid request", Data: data}
}

func ErrMethodNotFound(method string) *Error {
	return &Error{Code: CodeMethodNotFound, Message: "Method not found", Data: method}
}

func ErrInvalidParams(data any) *Error {
	return &Error{Code: CodeInvalidParams, Message: "Invalid params", Data: data}
}

func ErrInternalError(data any) *Error {
	return &Error{Code: CodeInternalError, Message: "Internal error", Data: data}
}

func ErrEvalNotFound(path string) *Error {
	return &Error{Code: CodeEvalNotFound, Message: "Eval not found", Data: path}
}

func ErrValidationFailed(data any) *Error {
	return &Error{Code: CodeValidationFailed, Message: "Validation failed", Data: data}
}

func ErrRunFailed(data any) *Error {
	return &Error{Code: CodeRunFailed, Message: "Run failed", Data: data}
}
