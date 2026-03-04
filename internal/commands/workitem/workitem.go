package workitem

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed workitem.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "workitem",
		Aliases: []string{"wi", "work-item"},
		Short:   "List and view Azure DevOps work items",
		Long:    doc,
	}
	cmd.AddCommand(
		listCmd(),
		viewCmd(),
		createCmd(),
	)
	return cmd
}
