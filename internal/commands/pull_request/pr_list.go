package pull_request

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

var prList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pull requests in the repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return List(cmd.Context())
	},
}

func List(ctx context.Context) error {
	cfg := config.From(ctx)
	list, err := rest.New(cfg.Tenant).
		Git().
		PRs(cfg.Repository).
		List(ctx, git_prs.ListQuery{
			SearchCriteria: &git_prs.SearchCriteria{
				Status: util.Ptr(models.PullRequestStatusActive),
			},
		})
	if err != nil {
		log.Error(err)
		return err
	}

	return util.DumpJSON(list)
}
