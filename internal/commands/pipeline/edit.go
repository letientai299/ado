package pipeline

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/spf13/cobra"
)

//go:embed edit.md
var editDoc string

type EditConfig struct {
	filterConfig
}

func editCmd() *cobra.Command {
	opts := &EditConfig{}

	cmd := &cobra.Command{
		Use:     "edit <id|name>",
		Aliases: []string{"e"},
		Short:   "Open pipeline YAML definition in editor",
		Long:    editDoc,
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}

			return newEditProcessor(c).process(args)
		},
	}
	return cmd
}

func newEditProcessor(c *common[*EditConfig]) *editProcessor {
	return &editProcessor{common: c}
}

type editProcessor struct {
	*common[*EditConfig]
}

func (e editProcessor) process(args []string) error {
	// Try if the first arg is a pipeline ID
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			m, err := e.client.Pipelines().Definitions(e.cfg.Repository).ByID(e.ctx, int32(id))
			if err == nil {
				return e.editOne(*m)
			}
		}
	}

	// Fallback to list/filter logic
	pipelines, err := e.list()
	if err != nil {
		return err
	}

	pipelines = e.filter(pipelines)

	switch len(pipelines) {
	case 0:
		return errors.New("no pipeline found matching the criteria")
	case 1:
		return e.editByID(pipelines[0].Id)
	default:
		return e.pick(pipelines)
	}
}

func (e editProcessor) pick(pipelines []Pipeline) error {
	selected := ui.Pick(pipelines, ui.PickConfig[Pipeline]{
		Title: "Select a pipeline to edit",
		Render: func(w io.Writer, p Pipeline, matches []int) {
			p.Name = styles.HighlightMatch(p.Name, matches)
			util.PanicIf(styles.Render(w, pipelinePickTpl, p))
		},
		FilterValue: func(p Pipeline) string { return strings.ToLower(p.Name) },
	})

	if selected.IsSome() {
		p := selected.Get()
		return e.editByID(p.Id)
	}
	return errors.New("no pipeline selected")
}

func (e editProcessor) editByID(id int32) error {
	m, err := e.client.Pipelines().Definitions(e.cfg.Repository).ByID(e.ctx, id)
	if err != nil {
		return err
	}

	return e.editOne(*m)
}

func (e editProcessor) editOne(m models.BuildDefinition) error {
	if m.Process == nil || m.Process.YamlFilename == "" {
		return errors.New("pipeline does not have a YAML definition")
	}

	yamlPath := m.Process.YamlFilename

	// Remove the leading slash if present
	yamlPath = strings.TrimPrefix(yamlPath, "/")

	// Find the git root directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Try to find the file relative to the current directory
	fullPath := filepath.Join(cwd, yamlPath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("YAML file not found: %s", yamlPath)
	}

	fmt.Printf("Opening %s\n", yamlPath)
	return editor.Open(e.cfg.Editor, fullPath)
}
