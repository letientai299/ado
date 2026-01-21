package pipeline

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/fp"
	"github.com/spf13/cobra"
)

const (
	outputJSON   = "json"
	outputYAML   = "yaml"
	outputSimple = "simple"
)

//go:embed list_simple.tpl
var listSimpleTpl string

//go:embed list.md
var listDoc string

// ListConfig holds configuration for the pipeline list command.
type ListConfig struct {
	DefaultOutput         string            `yaml:"default_output"          json:"default_output"`
	CustomOutputTemplates map[string]string `yaml:"custom_output_templates" json:"custom_output_templates"`

	filterConfig `yaml:"-"`
	output       *util.EnumFlag[string] `yaml:"-"`
}

func (l *ListConfig) OnResolved(c *cobra.Command) error {
	for name := range l.CustomOutputTemplates {
		l.output.AddAllowed(name)
	}

	if l.DefaultOutput != "" {
		flag := c.PersistentFlags().Lookup("output")
		if flag != nil {
			flag.DefValue = l.DefaultOutput
		}
		if !c.Flags().Changed("output") {
			_ = l.output.Set(l.DefaultOutput)
		}
	}

	return l.output.Validate()
}

func listCmd() *cobra.Command {
	opts := defaultListConfig()

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List pipelines in the repo",
		Long:    listDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return listProcessor{c}.process()
		},
	}

	flags := cmd.PersistentFlags()
	flags.VarP(opts.output, "output", "o", "output format")
	opts.output.RegisterCompletion(cmd, "output")
	return cmd
}

func defaultListConfig() *ListConfig {
	opts := &ListConfig{
		DefaultOutput:         outputSimple,
		CustomOutputTemplates: make(map[string]string),
		output: util.NewEnumFlag(outputSimple, outputJSON, outputYAML).
			Default(outputSimple),
	}

	config.Register(config.CommandConfig{
		Path:   "pipeline.list",
		Target: opts,
	})

	return opts
}

type listProcessor struct {
	*common[*ListConfig]
}

func (l listProcessor) process() error {
	pipelines, err := l.list()
	if err != nil {
		return err
	}

	pipelines = l.filter(pipelines)
	return l.render(pipelines)
}

func (l listProcessor) render(raw []models.BuildDefinition) error {
	all := fp.Map(raw, func(i models.BuildDefinition) Pipeline {
		return toPipeline(i, l.baseURL)
	})
	output := strings.ToLower(l.opts.output.Value())
	switch output {
	case outputYAML:
		return styles.DumpYAML(all)
	case outputJSON:
		return styles.DumpJSON(all)
	case outputSimple:
		return l.renderTemplate(listSimpleTpl, all)
	default:
		if tpl, ok := l.opts.CustomOutputTemplates[output]; ok {
			return l.renderTemplate(tpl, all)
		}
	}

	return util.StrErr("unknown output format: " + l.opts.output.Value())
}

func (l listProcessor) renderTemplate(tpl string, all []Pipeline) error {
	return styles.RenderOut(tpl, all, template.FuncMap{})
}
