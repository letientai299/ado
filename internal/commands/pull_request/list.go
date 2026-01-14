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
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

const (
	outputJSON   = "json"
	outputSimple = "simple"
)

func ListCmd() *cobra.Command {
	opt := listOptions{
		output: outputSimple,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List pull requests in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return List(cmd.Context(), opt)
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

func List(ctx context.Context, opts listOptions) error {
	log.Debugf("listing options: %+v", opts)

	cfg := config.From(ctx)
	criteria := &git_prs.SearchCriteria{
		Status: util.Ptr(models.PullRequestStatusActive),
	}

	all, err := rest.New(cfg.Tenant).
		Git().
		PRs(cfg.Repository).
		List(ctx, git_prs.ListQuery{SearchCriteria: criteria})
	if err != nil {
		log.Error(err)
		return err
	}

	opts.username = cfg.Username
	all = slices.DeleteFunc(all, opts.match)
	return render(ctx, all, opts.output)
}

func render(ctx context.Context, all []models.GitPullRequest, output string) error {
	if output == outputJSON {
		return util.DumpJSON(all)
	}

	if output == outputSimple {
		return renderSimple(ctx, all)
	}

	return util.StrErr("unknown output format: " + output)
}

func renderSimple(ctx context.Context, all []models.GitPullRequest) error {
	cfg := config.From(ctx)
	baseURL, _ := url.JoinPath(cfg.Repository.WebURL(), "pullRequest")
	for _, pr := range all {
		if pr.IsDraft {
			fmt.Print("DRAFT | ")
		}
		fmt.Println(pr.Title)
		fmt.Println("  " + pr.CreatedBy.DisplayName)
		fmt.Println("  " + baseURL + "/" + strconv.Itoa(pr.PullRequestId))
		fmt.Println("  " + pr.Url)
	}
	return nil
}

type listOptions struct {
	filterOptions
	output string
}

type filterOptions struct {
	mine     bool
	username string
	draft    bool
}

func (f filterOptions) match(pr models.GitPullRequest) bool {
	if !f.draft && pr.IsDraft {
		return true
	}

	if f.mine {
		// NOTE (tai): the UniqueName might not be the az account username in some other ADO org setup.
		return pr.CreatedBy.UniqueName != f.username
	}

	return false
}
