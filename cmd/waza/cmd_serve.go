package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spboyer/waza/internal/jsonrpc"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	var tcpAddr string
	var tcpAllowRemote bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a JSON-RPC 2.0 server for IDE integration",
		Long: `Start a JSON-RPC 2.0 server for IDE integration.

By default, the server communicates over stdin/stdout using newline-delimited JSON.
This enables VS Code, JetBrains, and other editors to run evals programmatically.

Use --tcp to start a TCP server instead (useful for debugging).
TCP defaults to loopback (127.0.0.1) for security. Use --tcp-allow-remote to bind
to all interfaces.

Supported methods:
  eval.run       Run an eval (returns run ID)
  eval.list      List available evals in a directory
  eval.get       Get eval details
  eval.validate  Validate an eval spec
  task.list      List tasks for an eval
  task.get       Get task details
  run.status     Get run status
  run.cancel     Cancel a running eval`,
		RunE: func(cmd *cobra.Command, args []string) error {
			registry := jsonrpc.NewMethodRegistry()
			hctx := jsonrpc.NewHandlerContext()
			jsonrpc.RegisterHandlers(registry, hctx)

			logger := slog.Default()
			server := jsonrpc.NewServer(registry, logger)

			if tcpAddr != "" {
				tcpAddr = resolveTCPAddr(tcpAddr, tcpAllowRemote, logger)

				listener, err := jsonrpc.NewTCPListener(tcpAddr, server)
				if err != nil {
					return fmt.Errorf("failed to start TCP server: %w", err)
				}
				defer listener.Close() //nolint:errcheck
				fmt.Fprintf(os.Stderr, "JSON-RPC server listening on %s\n", listener.Addr())
				return listener.Serve()
			}

			fmt.Fprintln(os.Stderr, "JSON-RPC server running on stdio")
			server.ServeStdio(os.Stdin, os.Stdout)
			return nil
		},
	}

	cmd.Flags().StringVar(&tcpAddr, "tcp", "", "TCP address to listen on (e.g., :9000)")
	cmd.Flags().BoolVar(&tcpAllowRemote, "tcp-allow-remote", false,
		"Allow binding to non-loopback addresses (WARNING: exposes the server to the network with no authentication)")

	return cmd
}

// resolveTCPAddr ensures TCP addresses default to loopback unless --tcp-allow-remote is set.
func resolveTCPAddr(addr string, allowRemote bool, logger *slog.Logger) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// Likely just a port like "9000"; treat as ":9000".
		host = ""
		port = addr
	}

	if allowRemote {
		logger.Warn("TCP server binding to all interfaces — no authentication is provided",
			"address", addr)
		return addr
	}

	// Default to loopback if no host specified or if 0.0.0.0/:: is used without --tcp-allow-remote.
	if host == "" || host == "0.0.0.0" || host == "::" {
		logger.Info("JSON-RPC server listening on TCP (local only)")
		return net.JoinHostPort("127.0.0.1", port)
	}

	return addr
}
