package jsonrpc

import (
	"context"
	"encoding/json"
)

// Handler processes a JSON-RPC request and returns a result or error.
type Handler func(ctx context.Context, params json.RawMessage) (any, *Error)

// MethodRegistry maps method names to handlers.
type MethodRegistry struct {
	methods map[string]Handler
}

// NewMethodRegistry creates an empty registry.
func NewMethodRegistry() *MethodRegistry {
	return &MethodRegistry{methods: make(map[string]Handler)}
}

// Register adds a handler for a method name.
func (r *MethodRegistry) Register(method string, handler Handler) {
	r.methods[method] = handler
}

// Lookup returns the handler for a method, or nil if not found.
func (r *MethodRegistry) Lookup(method string) Handler {
	return r.methods[method]
}

// Methods returns all registered method names.
func (r *MethodRegistry) Methods() []string {
	names := make([]string, 0, len(r.methods))
	for name := range r.methods {
		names = append(names, name)
	}
	return names
}
