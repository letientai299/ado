package workitem

import (
	_ "embed"
	"fmt"
	"strconv"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/spf13/cobra"
)

//go:embed delete.md
var deleteDoc string

type DeleteConfig struct {
	destroy bool
	yes     bool
}

func deleteCmd() *cobra.Command {
	opts := &DeleteConfig{}

	cmd := &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a work item",
		Long:    deleteDoc,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newDeleteProcessor(c).process(args[0])
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.destroy, "destroy", false, "permanently delete (cannot be recovered)")
	flags.BoolVarP(&opts.yes, "yes", "y", false, "skip confirmation prompt")

	return cmd
}

type deleteProcessor struct {
	*common[*DeleteConfig]
}

func newDeleteProcessor(c *common[*DeleteConfig]) *deleteProcessor {
	return &deleteProcessor{common: c}
}

func (p *deleteProcessor) process(idArg string) error {
	id, err := strconv.Atoi(idArg)
	if err != nil {
		return fmt.Errorf("invalid work item ID: %q", idArg)
	}

	// Fetch the work item first to show its title in the confirmation
	wi, err := p.client.WorkItems(p.cfg.Repository).ByID(p.ctx, id, models.WorkItemExpandNone)
	if err != nil {
		return fmt.Errorf("work item #%d not found: %w", id, err)
	}

	title := getStringField(*wi, models.FieldTitle)
	wiType := getStringField(*wi, models.FieldWorkItemType)

	if !p.opts.yes {
		action := "Delete"
		if p.opts.destroy {
			action = styles.Error("Permanently destroy")
		}
		msg := fmt.Sprintf("%s %s #%d %q?", action, wiType, id, title)
		if !ui.Confirm(msg, false) {
			return nil
		}
	}

	resp, err := p.client.WorkItems(p.cfg.Repository).Delete(p.ctx, id, p.opts.destroy)
	if err != nil {
		return fmt.Errorf("failed to delete work item #%d: %w", id, err)
	}

	if p.opts.destroy {
		fmt.Printf("Permanently destroyed #%d %s\n", resp.ID, styles.H1(title))
	} else {
		fmt.Printf("Deleted #%d %s (moved to Recycle Bin)\n", resp.ID, styles.H1(title))
	}

	return nil
}
