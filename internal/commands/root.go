package commands

import (
	_ "embed"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/letientai299/ado/internal/commands/api"
	"github.com/letientai299/ado/internal/commands/config_cmd"
	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/commands/workitem"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/profiling"
	"github.com/spf13/cobra"
)

var (
	//go:embed root.md
	doc string
	//go:embed templates.md
	templatesDoc string
	//go:embed usage.tpl
	usageTemplate string
	//go:embed help.tpl
	helpTemplate string
)

const groupExperimental = "Experimental"

var nonAdoCommands = []string{"completion", "version"}

var initOnce sync.Once

func Root() *cobra.Command {
	var stopProfiling profiling.StopFn

	cmdName := filepath.Base(os.Args[0])
	cmdName = cmdName[:len(cmdName)-len(filepath.Ext(cmdName))]

	root := &cobra.Command{
		Use:   cmdName,
		Short: "Azure DevOps CLI",
		Long:  doc,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			c := cmd
			for c.HasParent() {
				if slices.Contains(nonAdoCommands, c.Name()) {
					return nil
				}
				c = c.Parent()
			}
			stopProfiling = profiling.Start(cmd)
			return initConfig(cmd)
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			if stopProfiling != nil {
				stopProfiling()
			}
		},
		SilenceUsage: true,
	}

	addTemplateHelpers()
	root.SetHelpTemplate(helpTemplate)
	root.SetUsageTemplate(usageTemplate)
	root.SetHelpFunc(helpFunc(root.HelpFunc()))

	root.AddCommand(
		config_cmd.Cmd(),
		pull_request.Cmd(),
		Version(),
		util.HelpTopic("templates", templatesDoc),
	)

	addGroup(root, groupExperimental,
		api.Cmd(),
		pipeline.Cmd(),
		workitem.Cmd(),
	)

	config.AddGlobalFlags(root)
	profiling.RegisterFlag(root)
	return root
}

func addGroup(root *cobra.Command, groupId string, cs ...*cobra.Command) {
	root.AddGroup(&cobra.Group{
		ID:    groupId,
		Title: groupId,
	})

	for _, c := range cs {
		root.AddCommand(c)
		c.GroupID = groupId
	}
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
