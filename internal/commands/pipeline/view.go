package pipeline

import (
	_ "embed"
	"errors"
	"fmt"
	"strconv"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/util/sh"
	"github.com/spf13/cobra"
)

//go:embed view.md
var viewDoc string

type ViewConfig struct {
	filterConfig
}

func viewCmd() *cobra.Command {
	opts := &ViewConfig{}

	cmd := &cobra.Command{
		Use:     "view <id|name>",
		Aliases: []string{"v", "browse"},
		Short:   "Open pipeline build page in browser",
		Long:    viewDoc,
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}

			return newViewProcessor(c).process(args)
		},
	}
	return cmd
}

func newViewProcessor(c *common[*ViewConfig]) *viewProcessor {
	return &viewProcessor{common: c}
}

type viewProcessor struct {
	*common[*ViewConfig]
}

func (v viewProcessor) process(args []string) error {
	// Try if the first arg is a pipeline ID
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			m, err := v.client.Pipelines().Definitions(v.cfg.Repository).ByID(v.ctx, int32(id))
			if err == nil {
				return v.openOne(*m)
			}
		}
	}

	// Fallback to list/filter logic
	pipelines, err := v.list()
	if err != nil {
		return err
	}

	pipelines = v.filter(pipelines)

	switch len(pipelines) {
	case 0:
		return errors.New("no pipeline found matching the criteria")
	case 1:
		return v.openByID(pipelines[0].Id)
	default:
		return v.pickView(pipelines)
	}
}

func (v viewProcessor) pickView(pipelines []models.BuildDefinition) error {
	selected := pick(pipelines)
	if selected.IsSome() {
		p := selected.Get()
		return v.openByID(p.Id)
	}
	return errors.New("no pipeline selected")
}

func (v viewProcessor) openByID(id int32) error {
	m, err := v.client.Pipelines().Definitions(v.cfg.Repository).ByID(v.ctx, id)
	if err != nil {
		return err
	}

	return v.openOne(*m)
}

func (v viewProcessor) openOne(m models.BuildDefinition) error {
	p := toPipeline(m, v.baseURL)
	fmt.Println(p.WebURL)
	return sh.Browse(p.WebURL)
}
