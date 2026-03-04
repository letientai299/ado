package workitem

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed create.md
var createDoc string

type CreateConfig struct {
	title     string
	desc      string
	wiType    string
	assignee  string
	area      string
	iteration string
	yes       bool
	browse    bool
}

func createCmd() *cobra.Command {
	opts := &CreateConfig{}

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c", "new"},
		Short:   "Create a new work item",
		Long:    createDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newCreateProcessor(c).process()
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.title, "title", "", "work item title")
	flags.StringVarP(&opts.desc, "description", "d", "", "work item description")
	flags.StringVarP(
		&opts.wiType,
		"type",
		"t",
		"Task",
		"work item type (e.g., Bug, Task, \"User Story\")",
	)
	flags.StringVarP(&opts.assignee, "assignee", "A", "", "assign to user (display name or email)")
	flags.StringVar(&opts.area, "area", "", "area path (e.g., Project\\Team)")
	flags.StringVar(&opts.iteration, "iteration", "", "iteration path (e.g., Project\\Sprint 1)")
	flags.BoolVarP(&opts.yes, "yes", "y", false, "skip confirmation prompt")
	flags.BoolVarP(&opts.browse, "browse", "b", false, "open work item in browser after creating")

	return cmd
}

type createProcessor struct {
	*common[*CreateConfig]
}

func newCreateProcessor(c *common[*CreateConfig]) *createProcessor {
	return &createProcessor{common: c}
}

func (p *createProcessor) process() error {
	fmt.Printf("Return here")
	return nil
}
