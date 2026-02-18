package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spboyer/waza/internal/jsonrpc"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	var tcpAddr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a JSON-RPC 2.0 server for IDE integration",
		Long: `Start a JSON-RPC 2.0 server for IDE integration.

By default, the server communicates over stdin/stdout using newline-delimited JSON.
This enables VS Code, JetBrains, and other editors to run evals programmatically.

Use --tcp to start a TCP server instead (useful for debugging and remote access).

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
				listener, err := jsonrpc.NewTCPListener(tcpAddr, server)
				if err != nil {
					return fmt.Errorf("failed to start TCP server: %w", err)
				}
				defer listener.Close()
				fmt.Fprintf(os.Stderr, "JSON-RPC server listening on %s\n", listener.Addr())
				return listener.Serve()
			}

			fmt.Fprintln(os.Stderr, "JSON-RPC server running on stdio")
			server.ServeStdio(os.Stdin, os.Stdout)
			return nil
		},
	}

	cmd.Flags().StringVar(&tcpAddr, "tcp", "", "TCP address to listen on (e.g., :9000)")

	return cmd
}
