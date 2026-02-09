package tokens

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Token management for markdown files",
		Long: `Analyze token counts in markdown files. Subcommands:
  count     Count tokens in markdown files`,
	}
	cmd.AddCommand(newCountCmd())
	return cmd
}
