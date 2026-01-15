package pull_request

import (
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a pull request",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return c
}
