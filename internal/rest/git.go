package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
	"github.com/letientai299/ado/internal/rest/git_prs"
)

// Git is https://learn.microsoft.com/en-us/rest/api/azure/devops/git
type Git struct {
	client Client
}

func (g Git) PRs(repo config.Repository) GitPRs {
	baseUrl, _ := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/git/repositories",
		repo.Repo,
		"pullrequests",
	)

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
	return httpGet[models.GitPullRequest](
		ctx,
		g.client,
		prURL,
		_shared.BoolQ("includeWorkItemRefs"),
		_shared.BoolQ("includeCommits"),
	)
}

// List call
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests
func (g GitPRs) List(
	ctx context.Context,
	q git_prs.ListQuery,
) ([]models.GitPullRequest, error) {
	list, err := httpGet[List[models.GitPullRequest]](ctx, g.client, g.baseUrl, q)
	if err != nil {
		return nil, err
	}
	return list.Value, err
}

type List[T any] struct {
	Value []T `json:"value"`
	Count int `json:"count"`
}
