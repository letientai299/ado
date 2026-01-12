package pull_request

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "pull-request",
	Aliases: []string{"pr"},
	Short:   "List, view, create or manipulate pull requests",
}

var prBrowse = &cobra.Command{
	Use:     "browse",
	Aliases: []string{"open", "o"},
	Short:   "Browse a pull request in the web",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var prUpdate = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update a pull request",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	Cmd.AddCommand(prList, prCreate, prUpdate, prBrowse)
}
