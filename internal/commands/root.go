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
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed root.md
var doc string

//go:embed usage.tpl
var usageTemplate string

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               os.Args[0],
		Short:             "Azure DevOps CLI",
		Long:              doc,
		PersistentPreRunE: config.Resolve,
	}

	rootCmd.SetUsageTemplate(usageTemplate)
	rootCmd.SetHelpFunc(helpFunc)

	rootCmd.AddCommand(
		pull_request.Cmd(),
		pipeline.Cmd(),
		Doctor(),
	)

	config.AddGlobalFlags(rootCmd)
	return rootCmd
}

func helpFunc(cmd *cobra.Command, _ []string) {
	addTemplateHelpers()
	help := cmd.Long
	if help == "" {
		help = cmd.Short
	}

	rendered, _ := styles.Markdown(help)
	cmd.Print(strings.TrimLeftFunc(rendered, unicode.IsSpace))
	cmd.Println(cmd.UsageString())
}

func addTemplateHelpers() {
	cobra.AddTemplateFunc("headingStyle", styles.HeadingStyle)
	cobra.AddTemplateFunc("flagStyle", styles.FlagStyle)
	cobra.AddTemplateFunc("cmdStyle", styles.CmdStyle)
	cobra.AddTemplateFunc("flags", flagSlice)
	cobra.AddTemplateFunc("flagName", flagName)
	cobra.AddTemplateFunc("indent", util.Indent)
	cobra.AddTemplateFunc("wrap", styles.Wrap)
}

func flagSlice(fs *pflag.FlagSet) []*pflag.Flag {
	var list []*pflag.Flag
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			list = append(list, f)
		}
	})
	return list
}

func flagName(f *pflag.Flag) string {
	var s string
	if f.Shorthand != "" {
		s = "-" + f.Shorthand + ", --" + f.Name
	} else {
		s = "    --" + f.Name
	}
	if f.Value.Type() != "bool" {
		s += " " + f.Value.Type()
	}
	return s
}
