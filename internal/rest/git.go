package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
)

// Git is https://learn.microsoft.com/en-us/rest/api/azure/devops/git
type Git struct {
	client Client
}

func (g Git) PRs(repo config.Repository) GitPRs {
	baseUrl, err := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/git/repositories",
		repo.Repo,
		"pullrequests",
	)
	if err != nil {
		log.Fatal(err)
	}

	return GitPRs{client: g.client, baseUrl: baseUrl}
}

// GitPRs is https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests
type GitPRs struct {
	client  Client
	baseUrl string
}

// ByID call
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-request
// with some default query params.
func (g GitPRs) ByID(ctx context.Context, prID int32) (*models.GitPullRequest, error) {
	prURL, _ := url.JoinPath(g.baseUrl, strconv.FormatInt(int64(prID), 10))
	return httpGet[models.GitPullRequest](ctx, g.client, prURL,
		// query params
		"includeWorkItemRefs", true,
		"includeCommits", true,
	)
}
