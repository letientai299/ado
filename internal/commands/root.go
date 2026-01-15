package commands

import (
	_ "embed"
	"os"
	"sync"

	"github.com/letientai299/ado/internal/commands/config_cmd"
	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/config"
	"github.com/spf13/cobra"
)

var (
	//go:embed root.md
	doc string
	//go:embed usage.tpl
	usageTemplate string
	//go:embed help.tpl
	helpTemplate string
)

var initOnce sync.Once

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   os.Args[0],
		Short: "Azure DevOps CLI",
		Long:  doc,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return initConfig(cmd.Root())
		},
		SilenceUsage: true,
	}

	root.SetHelpTemplate(helpTemplate)
	root.SetUsageTemplate(usageTemplate)
	root.SetHelpFunc(helpFunc(root.HelpFunc()))

	root.AddCommand(
		pull_request.Cmd(),
		pipeline.Cmd(),
		config_cmd.Cmd(),
		Doctor(),
		Version(),
	)

	config.AddGlobalFlags(root)
	return root
}

func initConfig(cmd *cobra.Command) error {
	var err error
	initOnce.Do(func() {
		if err = config.Resolve(cmd, nil); err == nil {
			addTemplateHelpers()
		}
	})
	return err
}
