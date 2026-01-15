package pull_request

import "github.com/spf13/cobra"

func updateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update a pull request",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	}
	return c
}
