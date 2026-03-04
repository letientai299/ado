package pipeline

import (
	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed pipeline.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipeline",
		Long:    doc,
		Short:   "list, view, edit pipelines",
		Aliases: []string{"pl"},
	}
	cmd.AddCommand(
		listCmd(),
		buildsCmd(),
		viewCmd(),
		editCmd(),
		logsCmd(),
	)
	return cmd
}
