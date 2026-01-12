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
		Short:   "list, view, run pipeline",
		Aliases: []string{"pp"},
	}
	cmd.AddCommand(ppList, ppRun, ppBrowse)
	return cmd
}

var ppList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list pull requests in the repo",
	RunE: func(cmd *cobra.Command, args []string) error {return nil },
}

var ppRun = &cobra.Command{
	Use:     "run",
	Aliases: []string{"c"},
	Short:   "create a pull request",
	RunE: func(cmd *cobra.Command, args []string) error {return nil },
}

var ppBrowse = &cobra.Command{
	Use:     "browse",
	Aliases: []string{"u"},
	Short:   "browse recent runs of a pipeline on the web",
	RunE: func(cmd *cobra.Command, args []string) error {return nil },
}
