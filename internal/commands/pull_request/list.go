package pull_request

import (
	"context"
	_ "embed"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

const (
	outputJSON   = "json"
	outputYAML   = "yaml"
	outputSimple = "simple"
)

//go:embed list_simple.tpl
var listSimpleTpl string

type PR = models.GitPullRequest

// ListConfig holds configuration for the pr list command.
// These values can be set in the config file under "pull-request.list".
type ListConfig struct {
	// Default output format to use if not specified.
	DefaultOutput string `yaml:"default_output" json:"default_output"`
	// Custom output templates is a map of output format names to their templates.
	CustomOutputTemplates map[string]string `yaml:"custom_output_templates" json:"custom_output_templates"`

	/* filtering */
	mine     bool     // shows only your PRs
	draft    bool     // whether to include draft PRs
	keywords []string // keywords to do filter PRs title and description

	/* rendering */
	output string // output format to use
}

func (l *ListConfig) OnResolved(c *cobra.Command) error {
	// TODO (tai): doesn't work correctly, as the flag.Changed() isn't checked.
	fs := c.Flags()
	if !fs.Changed("output") {
		l.output = l.DefaultOutput
	}
	return nil
}

func listCmd() *cobra.Command {
	opts := &ListConfig{
		DefaultOutput:         outputSimple,
		CustomOutputTemplates: make(map[string]string),
	}

	config.Register(config.CommandConfig{
		Path:   "pull-request.list",
		Target: opts,
	})

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List pull requests in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := config.From(ctx)
			token, err := cfg.Token()
			if err != nil {
				return err
			}
			client := rest.New(token)
			baseURL, _ := url.JoinPath(cfg.Repository.WebURL(), "pullRequest")
			opts.keywords = args
			return listProcessor{
				opts:    opts,
				cfg:     cfg,
				client:  client,
				baseURL: baseURL,
			}.process(ctx)
		},
	}

	flags := cmd.PersistentFlags()

	// filter flags
	flags.BoolVarP(&opts.mine, "mine", "m", false, "show only your PRs")
	flags.BoolVar(&opts.draft, "draft", false, "include draft PRs")

	// render flags
	flags.StringVarP(
		&opts.output,
		"output",
		"o",
		// TODO (tai): help should show resolved config value instead of hard-coded values
		opts.DefaultOutput,
		// TODO (tai): remove json format
		"output format (builtin: simple, json, yaml)",
	)
	return cmd
}

type listProcessor struct {
	opts    *ListConfig
	client  *rest.Client
	cfg     *config.Config
	baseURL string
}

func (l listProcessor) process(ctx context.Context) error {
	prs, err := l.query(ctx)
	if err != nil {
		return err
	}

	prs, err = l.filter(ctx, prs)
	if err != nil {
		return err
	}

	return l.render(prs)
}

func (l listProcessor) query(ctx context.Context) ([]models.GitPullRequest, error) {
	criteria := &git_prs.SearchCriteria{
		Status: util.Ptr(models.PullRequestStatusActive),
	}

	all, err := l.client.Git().
		PRs(l.cfg.Repository).
		List(ctx, git_prs.ListQuery{SearchCriteria: criteria})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return all, nil
}

func (l listProcessor) filter(ctx context.Context, all []PR) ([]PR, error) {
	var id *string
	if l.opts.mine {
		identity, err := l.client.Identity(ctx, l.cfg.Repository.Org)
		if err != nil {
			return nil, err
		}
		id = &identity.Id
	}

	return slices.DeleteFunc(all, func(pr PR) bool {
		if !l.opts.draft && pr.IsDraft {
			return true
		}

		if id != nil && pr.CreatedBy.Id != *id {
			return true
		}

		return !containsAll(pr, l.opts.keywords)
	}), nil
}

func containsAll(pr PR, keywords []string) bool {
	for _, pattern := range keywords {
		if !strings.Contains(pr.Title, pattern) && !strings.Contains(pr.Description, pattern) {
			return false
		}
	}
	return true
}

func (l listProcessor) render(all []PR) error {
	output := strings.ToLower(l.opts.output)
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

	return util.StrErr("unknown output format: " + l.opts.output)
}

func (l listProcessor) renderTemplate(tpl string, all []PR) error {
	return styles.RenderTemplate(tpl, all, template.FuncMap{
		"webURL": l.webURL,
	})
}

func (l listProcessor) webURL(pr PR) string {
	return l.baseURL + "/" + strconv.Itoa(pr.PullRequestId)
}
