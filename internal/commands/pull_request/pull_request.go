package pull_request

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed pull_request.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pr",
		Aliases: []string{"pull-request", "pull"},
		Short:   "List, view, create or manipulate pull requests",
		Long:    doc,
	}
	cmd.AddCommand(prList, prCreate, prUpdate, prView)
	return cmd
}

var prUpdate = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update a pull request",
	RunE:    func(cmd *cobra.Command, args []string) error { return nil },
}

func init() {
}
