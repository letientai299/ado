package pull_request

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed pull_request.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull-request",
		Aliases: []string{"pr", "pull"},
		Short:   "List, view, create or manipulate pull requests",
		Long:    doc,
	}
	cmd.AddCommand(
		listCmd(),
		viewCmd(),
		createCmd(),
		updateCmd(),
	)
	return cmd
}
