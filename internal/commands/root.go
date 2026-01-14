package commands

import (
	_ "embed"
	"os"

	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/config"
	"github.com/spf13/cobra"
)

//go:embed root.md
var doc string

//go:embed usage.tpl
var usageTemplate string

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:               os.Args[0],
		Short:             "Azure DevOps CLI",
		Long:              doc,
		PersistentPreRunE: config.Resolve,
		SilenceUsage:      true,
	}

	root.SetUsageTemplate(usageTemplate)
	root.SetHelpFunc(prettifyHelp(root.HelpFunc()))

	root.AddCommand(
		pull_request.Cmd(),
		pipeline.Cmd(),
		Doctor(),
	)

	config.AddGlobalFlags(root)
	return root
}
