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

// RepoInfo retrieves repository information including ID.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get-repository
func (g Git) RepoInfo(ctx context.Context, repo config.Repository) (*models.GitRepository, error) {
	repoURL, _ := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/git/repositories",
		repo.Name,
	)
	return httpGet[models.GitRepository](ctx, g.client, repoURL)
}

func (g Git) PRs(repo config.Repository) GitPRs {
	baseUrl, _ := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/git/repositories",
		repo.Name,
		"pullrequests",
	)

	return GitPRs{client: g.client, baseUrl: baseUrl, org: repo.Org}
}

// GitPRs is https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests
type GitPRs struct {
	client  Client
	baseUrl string
	org     string
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
		_shared.Bool("includeWorkItemRefs"),
		_shared.Bool("includeCommits"),
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

// Create call
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/create
func (g GitPRs) Create(
	ctx context.Context,
	pr models.GitPullRequest,
) (*models.GitPullRequest, error) {
	return httpPost[models.GitPullRequest](ctx, g.client, g.baseUrl, pr)
}

// Update call
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/update
func (g GitPRs) Update(
	ctx context.Context,
	id int32,
	pr models.GitPullRequest,
) (*models.GitPullRequest, error) {
	prURL, _ := url.JoinPath(g.baseUrl, strconv.FormatInt(int64(id), 10))
	return httpPatch[models.GitPullRequest](ctx, g.client, prURL, pr)
}

// Reviewers returns all reviewers for a PR, including those who voted optionally.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-reviewers/list
func (g GitPRs) Reviewers(ctx context.Context, prID int32) ([]models.IdentityRefWithVote, error) {
	reviewersURL, _ := url.JoinPath(
		g.baseUrl,
		strconv.FormatInt(int64(prID), 10),
		"reviewers",
	)
	list, err := httpGet[List[models.IdentityRefWithVote]](ctx, g.client, reviewersURL)
	if err != nil {
		return nil, err
	}
	return list.Value, nil
}

// Vote sets the current user's vote on a PR.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-reviewers/create-pull-request-reviewer
func (g GitPRs) Vote(
	ctx context.Context,
	prID int32,
	vote models.PrVote,
) (*models.IdentityRefWithVote, error) {
	identity, err := g.client.Identity(ctx, g.org)
	if err != nil {
		return nil, err
	}

	reviewerURL, _ := url.JoinPath(
		g.baseUrl,
		strconv.FormatInt(int64(prID), 10),
		"reviewers",
		identity.Id,
	)
	body := map[string]int{"vote": int(vote)}
	return httpPut[models.IdentityRefWithVote](ctx, g.client, reviewerURL, body)
}

type List[T any] struct {
	Value []T `json:"value"`
	Count int `json:"count"`
}
