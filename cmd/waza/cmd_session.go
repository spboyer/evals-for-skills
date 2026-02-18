package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spboyer/waza/internal/session"
	"github.com/spf13/cobra"
)

func newSessionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "View and manage session logs",
		Long: `View and manage session event logs.

Session logs are NDJSON files written during eval runs when --session-log is enabled.
They record the full lifecycle: session start, task execution, grader results, and completion.`,
	}

	cmd.AddCommand(newSessionListCommand())
	cmd.AddCommand(newSessionViewCommand())

	return cmd
}

func newSessionListCommand() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recorded session logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return err
			}

			files, err := session.ListSessions(absDir)
			if err != nil {
				return fmt.Errorf("listing sessions: %w", err)
			}

			if len(files) == 0 {
				fmt.Println("No session logs found.")
				return nil
			}

			fmt.Printf("%-40s %-8s %s\n", "File", "Events", "Modified")
			fmt.Println("─────────────────────────────────────────────────────────────────")
			for _, f := range files {
				fmt.Printf("%-40s %-8d %s\n", f.Name, f.NumEvents, f.ModTime.Format("2006-01-02 15:04:05"))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", ".", "Directory to search for session logs")

	return cmd
}

func newSessionViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <session-file>",
		Short: "View a session timeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			events, err := session.ReadEvents(path)
			if err != nil {
				return fmt.Errorf("reading session: %w", err)
			}

			session.RenderTimeline(os.Stdout, events)
			return nil
		},
	}

	return cmd
}
