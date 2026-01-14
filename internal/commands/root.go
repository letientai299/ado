package commands

import (
	_ "embed"
	"os"
	"strings"
	"unicode"

	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/styles"
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
		SilenceUsage: true,
	}

	root.SetUsageTemplate(usageTemplate)
	root.SetHelpFunc(prettifyHelp(root.HelpFunc()))

	// ensure config is resolved
	usageFunc := root.UsageFunc()
	root.SetUsageFunc(func(cmd *cobra.Command) error {
		if err := config.Resolve(cmd, nil); err != nil {
			return err
		}
		return usageFunc(cmd)
	})

	root.AddCommand(
		pull_request.Cmd(),
		pipeline.Cmd(),
		Doctor(),
	)

	config.AddGlobalFlags(root)
	return root
}

type helpFunc func(cmd *cobra.Command, args []string)

func prettifyHelp(defaultFn helpFunc) helpFunc {
	return func(cmd *cobra.Command, args []string) {
		if err := config.Resolve(cmd, args); err != nil {
			defaultFn(cmd, args)
			return
		}

		help := cmd.Long
		if help == "" {
			help = cmd.Short
		}

		rendered, _ := styles.Markdown(help)
		cmd.Print(strings.TrimLeftFunc(rendered, unicode.IsSpace))
		cmd.Println(cmd.UsageString())
	}
}
