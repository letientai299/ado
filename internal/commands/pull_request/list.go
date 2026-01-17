package pull_request

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

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

type PR = models.GitPullRequest

// listConfig holds configuration for the pr list command.
// These values can be set in the config file under "pull-request.list".
type listConfig struct {
	DefaultOutput         string            `yaml:"default_output"`
	CustomOutputTemplates map[string]string `yaml:"custom_output_templates"`

	output string // output format to use
	mine   bool   // shows only your PRs
	draft  bool   // whether to include draft PRs
}

func (l *listConfig) OnResolved(c *cobra.Command) error {
	// TODO (tai): doesn't work correctly, as the flag.Changed() isn't checked.
	fs := c.Flags()
	if !fs.Changed("output") {
		l.output = l.DefaultOutput
	}
	return nil
}

func listCmd() *cobra.Command {
	opts := &listConfig{
		DefaultOutput:         outputSimple,
		CustomOutputTemplates: make(map[string]string),
	}

	config.Register(config.CommandConfig{
		Path:   "pull-request.list",
		Target: opts,
	})

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List pull requests in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := config.From(ctx)
			client := rest.New(cfg.Token)
			return listProcessor{opts: opts, cfg: cfg, client: client}.process(ctx)
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
	opts   *listConfig
	client *rest.Client
	cfg    *config.Config
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

	return l.render(ctx, prs)
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

		return id != nil && pr.CreatedBy.Id != *id
	}), nil
}

func (l listProcessor) render(ctx context.Context, all []PR) error {
	switch strings.ToLower(l.opts.output) {
	case outputYAML:
		return styles.DumpYAML(all)
	case outputJSON:
		return styles.DumpJSON(all)
	case outputSimple:
		return renderSimple(ctx, all)
	default:
		return util.StrErr("unknown output format: " + l.opts.output)
	}
}

func renderSimple(ctx context.Context, all []PR) error {
	cfg := config.From(ctx)
	baseURL, _ := url.JoinPath(cfg.Repository.WebURL(), "pullRequest")
	for _, pr := range all {
		if pr.IsDraft {
			fmt.Print("DRAFT | ")
		}
		fmt.Println(pr.Title)
		fmt.Println("  " + pr.CreatedBy.DisplayName)
		fmt.Println("  " + baseURL + "/" + strconv.Itoa(pr.PullRequestId))
	}
	return nil
}
