package workitem

import (
	_ "embed"
	"fmt"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/styles"
	// "github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/letientai299/ado/internal/util/sh"
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
	opts := p.opts

	// If title is not provided, open editor to fill in all fields
	if opts.title == "" {
		if err := p.editFields(); err != nil {
			return err
		}
	}

	opts.title = strings.TrimSpace(opts.title)
	if opts.title == "" {
		return fmt.Errorf("title is required, use --title or fill it in the editor")
	}

	// if !opts.yes {
	// 	summary := fmt.Sprintf("Create %s: %q?", opts.wiType, opts.title)
	// 	if !ui.Confirm(summary, true) {
	// 		return nil
	// 	}
	// }
	fmt.Printf("Return here")
	return nil

	fields := []rest.JsonPatchOp{
		{Op: "add", Path: "/fields/" + models.FieldTitle, Value: opts.title},
	}

	for _, f := range []struct {
		field string
		value string
	}{
		{models.FieldDescription, opts.desc},
		{models.FieldAssignedTo, opts.assignee},
		{models.FieldAreaPath, opts.area},
		{models.FieldIterationPath, opts.iteration},
	} {
		if f.value != "" {
			fields = append(fields, rest.JsonPatchOp{
				Op: "add", Path: "/fields/" + f.field, Value: f.value,
			})
		}
	}

	wi, err := p.client.WorkItems(p.cfg.Repository).Create(p.ctx, opts.wiType, fields)
	if err != nil {
		return fmt.Errorf("failed to create work item: %w", err)
	}

	wiURL := fmt.Sprintf("%s/%d", p.baseURL, wi.ID)
	fmt.Printf("#%d %s\n", wi.ID, styles.H1(getStringField(*wi, models.FieldTitle)))
	fmt.Println(wiURL)

	if opts.browse {
		return sh.Browse(wiURL)
	}

	return nil
}
