package pipeline

import (
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/spf13/cobra"
)

//go:embed builds.md
var buildsDoc string

//go:embed builds_simple.tpl
var buildsSimpleTpl string

type BuildsConfig struct {
	filterConfig
	pipelineID int32
	top        int
	yamlFile   string
}

func buildsCmd() *cobra.Command {
	opts := &BuildsConfig{}

	cmd := &cobra.Command{
		Use:     "builds [keywords...]",
		Aliases: []string{"build", "b"},
		Short:   "List builds for a pipeline",
		Long:    buildsDoc,
		Args:    cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return buildsProcessor{c}.process(args)
		},
	}

	flags := cmd.Flags()
	flags.Int32VarP(&opts.pipelineID, "pipeline-id", "p", 0, "pipeline definition ID")
	flags.IntVarP(&opts.top, "number", "n", 10, "number of builds to show")
	flags.StringVar(&opts.yamlFile, "yaml", "", "filter pipeline by YAML file path")

	return cmd
}

type buildsProcessor struct {
	*common[*BuildsConfig]
}

func (b buildsProcessor) process(args []string) error {
	pipelineID := b.opts.pipelineID

	// If no --pipeline-id, resolve from keywords or picker
	if pipelineID == 0 {
		pipeline, err := b.selectPipeline(args)
		if err != nil {
			return err
		}
		pipelineID = pipeline.Id
	}

	builds, err := b.client.Builds().ForProject(b.cfg.Repository).List(b.ctx, rest.BuildListOptions{
		DefinitionID: pipelineID,
		Top:          b.opts.top,
	})
	if err != nil {
		return fmt.Errorf("failed to list builds: %w", err)
	}

	if len(builds) == 0 {
		return errors.New("no builds found for this pipeline")
	}

	return b.render(builds)
}

func (b buildsProcessor) selectPipeline(args []string) (*models.BuildDefinition, error) {
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			m, err := b.client.Pipelines().Definitions(b.cfg.Repository).ByID(b.ctx, int32(id))
			if err == nil {
				return m, nil
			}
		}
	}

	pipelines, err := b.list()
	if err != nil {
		return nil, err
	}

	if b.opts.yamlFile != "" {
		want := strings.TrimPrefix(strings.ToLower(gitcli.ResolveRepoRelativePath(b.opts.yamlFile)), "/")
		var matched []models.BuildDefinition
		for _, p := range pipelines {
			if p.Process != nil {
				got := strings.TrimPrefix(strings.ToLower(p.Process.YamlFilename), "/")
				if got == want {
					matched = append(matched, p)
				}
			}
		}
		pipelines = matched
	}

	pipelines = b.filter(pipelines)

	switch len(pipelines) {
	case 0:
		return nil, errors.New("no pipeline found matching the criteria")
	case 1:
		return &pipelines[0], nil
	default:
		selected := pick(pipelines)
		if selected.IsSome() {
			p := selected.Get()
			return &p, nil
		}
		return nil, errors.New("no pipeline selected")
	}
}

// BuildItem is the DTO for build display in the template.
type BuildItem struct {
	Id       int32  `yaml:"id"        json:"id"`
	Number   string `yaml:"number"    json:"number"`
	Status   string `yaml:"status"    json:"status"`
	Result   string `yaml:"result"    json:"result"`
	Branch   string `yaml:"branch"    json:"branch"`
	Commit   string `yaml:"commit"    json:"commit,omitempty"`
	Reason   string `yaml:"reason"    json:"reason"`
	Duration string `yaml:"duration"  json:"duration,omitempty"`
}

func toBuildItem(b models.Build) BuildItem {
	item := BuildItem{
		Id:     b.Id,
		Number: b.BuildNumber,
		Status: string(b.Status),
		Result: string(b.Result),
		Branch: strings.TrimPrefix(b.SourceBranch, "refs/heads/"),
		Reason: string(b.Reason),
	}

	if b.SourceVersionMessage != "" {
		msg := b.SourceVersionMessage
		if idx := strings.Index(msg, "\n"); idx > 0 {
			msg = msg[:idx]
		}
		item.Commit = msg
	}

	if b.StartTime != nil && b.FinishTime != nil {
		item.Duration = b.FinishTime.Sub(*b.StartTime).Truncate(time.Second).String()
	}

	return item
}

func (b buildsProcessor) render(builds []models.Build) error {
	items := make([]BuildItem, len(builds))
	for i, build := range builds {
		items[i] = toBuildItem(build)
	}
	return styles.RenderOut(buildsSimpleTpl, items, template.FuncMap{})
}
