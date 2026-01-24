package pull_request

import (
	_ "embed"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/git_prs"
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

// ListConfig holds configuration for the pr list command.
// These values can be set in the config file under "pull-request.list".
type ListConfig struct {
	// Default output format to use if not specified.
	DefaultOutput string `yaml:"default_output" json:"default_output"`
	// Custom output templates is a map of output format names to their templates.
	CustomOutputTemplates map[string]string `yaml:"custom_output_templates" json:"custom_output_templates"`

	filterConfig `yaml:"-"`
	output       *util.EnumFlag[string] `yaml:"-"` // output format to use
}

func (l *ListConfig) OnResolved(c *cobra.Command) error {
	// Add custom output formats from config
	for name := range l.CustomOutputTemplates {
		l.output.AddAllowed(name)
	}

	// Update default value if configured
	if l.DefaultOutput != "" {
		flag := c.PersistentFlags().Lookup("output")
		if flag != nil {
			flag.DefValue = l.DefaultOutput
		}
		// Set the value if not explicitly changed by the user
		if !c.Flags().Changed("output") {
			_ = l.output.Set(l.DefaultOutput)
		}
	}

	// Validate after all allowed values have been added
	return l.output.Validate()
}

func listCmd() *cobra.Command {
	opts := defaultListConfig()

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List pull requests in the repo",
		Long:    listDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newListProcessor(c).process()
		},
	}

	opts.RegisterFlags(cmd)
	flags := cmd.PersistentFlags()

	// render flags
	flags.VarP(opts.output, "output", "o", "output format")
	opts.output.RegisterCompletion(cmd, "output")
	return cmd
}

func newListProcessor(c *common[*ListConfig]) listProcessor {
	return listProcessor{c}
}

func defaultListConfig() *ListConfig {
	opts := &ListConfig{
		DefaultOutput:         outputSimple,
		CustomOutputTemplates: make(map[string]string),
		output: util.NewEnumFlag(outputSimple, outputJSON, outputYAML).
			Default(outputSimple),
	}

	config.Register(config.CommandConfig{
		Path:   "pull-request.list",
		Target: opts,
	})

	return opts
}

type listProcessor struct {
	*common[*ListConfig]
}

func (l listProcessor) process() error {
	prs, err := l.find()
	if err != nil {
		return err
	}

	return l.render(prs)
}

func (l listProcessor) find() ([]models.GitPullRequest, error) {
	prs, err := l.query()
	if err != nil {
		return nil, err
	}

	return l.filter(prs)
}

func (l listProcessor) query() ([]models.GitPullRequest, error) {
	criteria := &git_prs.SearchCriteria{
		Status: util.Ptr(models.PullRequestStatusActive),
	}

	all, err := l.client.Git().
		PRs(l.cfg.Repository).
		List(l.ctx, git_prs.ListQuery{SearchCriteria: criteria})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return all, nil
}

func (l listProcessor) filter(all []models.GitPullRequest) ([]models.GitPullRequest, error) {
	var id *string
	if l.opts.mine {
		identity, err := l.client.Identity(l.ctx, l.cfg.Repository.Org)
		if err != nil {
			return nil, err
		}
		id = &identity.Id
	}

	// Pre-convert keywords to lowercase once
	var lowerKeywords []string
	if len(l.opts.keywords) > 0 {
		lowerKeywords = make([]string, len(l.opts.keywords))
		for i, kw := range l.opts.keywords {
			lowerKeywords[i] = strings.ToLower(kw)
		}
	}

	return slices.DeleteFunc(all, func(m models.GitPullRequest) bool {
		if !l.opts.draft && m.IsDraft {
			return true
		}

		if id != nil && m.CreatedBy.Id != *id {
			return true
		}

		return !l.containsAll(m, lowerKeywords)
	}), nil
}

func (l listProcessor) containsAll(pr models.GitPullRequest, lowerKeywords []string) bool {
	if len(lowerKeywords) == 0 {
		return true
	}

	// Convert title and desc to lowercase once
	title := strings.ToLower(pr.Title)
	desc := strings.ToLower(pr.Description)

	for _, pattern := range lowerKeywords {
		if !strings.Contains(title, pattern) && !strings.Contains(desc, pattern) {
			return false
		}
	}
	return true
}

func (l listProcessor) render(all []models.GitPullRequest) error {
	// Fetch policy evaluations for all PRs
	evaluations, err := l.fetchEvaluations(all)
	if err != nil {
		log.Warn("failed to fetch policy evaluations", "error", err)
		// Continue without evaluations rather than failing
		evaluations = nil
	}

	// Get first PR's repository for URL construction (all PRs should be from same repo)
	var repo *models.GitRepository
	if len(all) > 0 && all[0].Repository != nil {
		repo = all[0].Repository
	}

	prs := fp.Map(all, converter(l.baseURL, l.cfg.Repository.Org, repo, evaluations))
	log.Debug("found pull requests", "count", len(prs))

	output := strings.ToLower(l.opts.output.Value())
	switch output {
	case outputYAML:
		return styles.DumpYAML(all)
	case outputJSON:
		return styles.DumpJSON(all)
	case outputSimple:
		return l.renderTemplate(listSimpleTpl, prs)
	default:
		if tpl, ok := l.opts.CustomOutputTemplates[output]; ok {
			return l.renderTemplate(tpl, prs)
		}
	}

	return util.StrErr("unknown output format: " + l.opts.output.Value())
}

func (l listProcessor) fetchEvaluations(
	prs []models.GitPullRequest,
) (map[int32][]models.PolicyEvaluationRecord, error) {
	if len(prs) == 0 {
		return nil, nil
	}

	var prIDs []int32
	for _, pr := range prs {
		prIDs = append(prIDs, pr.PullRequestId)
	}

	return l.client.Policy().Evaluations(l.cfg.Repository).
		List(l.ctx, prIDs...)
}

func (l listProcessor) renderTemplate(tpl string, all []PR) error {
	return styles.RenderOut(tpl, all)
}
