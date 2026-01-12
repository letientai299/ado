package commands

import (
	"os"

	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/config"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               os.Args[0],
		Short:             "Azure DevOps CLI",
		PersistentPreRunE: config.Resolve,
	}

	rootCmd.AddCommand(
		pull_request.Cmd,
		pipeline.Cmd,
		Doctor(),
	)

	config.AddGlobalFlags(rootCmd)
	return rootCmd
}
