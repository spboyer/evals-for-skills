package jsonrpc

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
)

// Server handles JSON-RPC 2.0 requests over a Transport.
type Server struct {
	registry *MethodRegistry
	logger   *slog.Logger
}

// NewServer creates a JSON-RPC server with the given method registry.
func NewServer(registry *MethodRegistry, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	return &Server{registry: registry, logger: logger}
}

// ServeTransport reads requests from the transport and writes responses.
// It runs until the transport's reader returns io.EOF or a read error.
func (s *Server) ServeTransport(t *Transport) {
	ctx := context.Background()

	for {
		req, err := t.ReadRequest()
		if err != nil {
			if err == io.EOF {
				return
			}
			s.logger.Debug("read error", "error", err)
			resp := &Response{
				JSONRPC: "2.0",
				Error:   ErrParseError(err.Error()),
				ID:      json.RawMessage("null"),
			}
			if writeErr := t.WriteResponse(resp); writeErr != nil {
				s.logger.Debug("write error", "error", writeErr)
			}
			return
		}

		// Validate JSON-RPC version
		if req.JSONRPC != "2.0" {
			resp := &Response{
				JSONRPC: "2.0",
				Error:   ErrInvalidRequest("jsonrpc field must be \"2.0\""),
				ID:      req.ID,
			}
			if writeErr := t.WriteResponse(resp); writeErr != nil {
				s.logger.Debug("write error", "error", writeErr)
				return
			}
			continue
		}

		// Look up method
		handler := s.registry.Lookup(req.Method)
		if handler == nil {
			resp := &Response{
				JSONRPC: "2.0",
				Error:   ErrMethodNotFound(req.Method),
				ID:      req.ID,
			}
			if writeErr := t.WriteResponse(resp); writeErr != nil {
				s.logger.Debug("write error", "error", writeErr)
				return
			}
			continue
		}

		// Execute handler
		result, rpcErr := handler(ctx, req.Params)
		resp := &Response{
			JSONRPC: "2.0",
			ID:      req.ID,
		}
		if rpcErr != nil {
			resp.Error = rpcErr
		} else {
			resp.Result = result
		}

		if writeErr := t.WriteResponse(resp); writeErr != nil {
			s.logger.Debug("write error", "error", writeErr)
			return
		}
	}
}

// ServeStdio runs the server on stdin/stdout.
func (s *Server) ServeStdio(stdin io.Reader, stdout io.Writer) {
	transport := NewTransport(stdin, stdout)
	s.ServeTransport(transport)
}
