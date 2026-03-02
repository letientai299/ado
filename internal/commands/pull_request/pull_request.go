package pull_request

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed pr.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "List, view, create or manipulate pull requests",
		Long:  doc,
	}
	cmd.AddCommand(
		listCmd(),
		viewCmd(),
		createCmd(),
		updateCmd(),
		analysisCmd(),
	)
	return cmd
}
