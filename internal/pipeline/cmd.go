package pipeline

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pp"},
	Short:   "list, view, run pipeline",
}

var ppList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list pull requests in the repo",
}

var ppRun = &cobra.Command{
	Use:     "run",
	Aliases: []string{"c"},
	Short:   "create a pull request",
}

var ppBrowse = &cobra.Command{
	Use:     "browse",
	Aliases: []string{"u"},
	Short:   "browse recent runs of a pipeline on the web",
}

func init() {
	Cmd.AddCommand(ppList, ppRun, ppBrowse)
}
