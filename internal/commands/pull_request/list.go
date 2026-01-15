package pull_request

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"

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

type listOptions struct {
	filterOptions
	output string
}

type filterOptions struct {
	mine  bool
	draft bool
}

func listCmd() *cobra.Command {
	opt := listOptions{
		output: outputSimple,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List pull requests in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg := config.From(ctx)
			client := rest.New(cfg.Token)
			return listProcessor{opts: opt, cfg: cfg, client: client}.process(ctx)
		},
	}

	flags := cmd.PersistentFlags()

	// filter flags
	flags.BoolVarP(&opt.mine, "mine", "m", false, "show only PRs created by you")
	flags.BoolVar(&opt.draft, "draft", false, "include draft PRs")

	// render flags
	flags.StringVarP(&opt.output, "output", "o", opt.output, "include draft PRs")
	return cmd
}

type listProcessor struct {
	opts   listOptions
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

	f := l.opts.filterOptions
	return slices.DeleteFunc(all, func(pr PR) bool {
		if !f.draft && pr.IsDraft {
			return true
		}

		return id != nil && pr.CreatedBy.Id != *id
	}), nil
}

func (l listProcessor) render(ctx context.Context, all []PR) error {
	switch l.opts.output {
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
