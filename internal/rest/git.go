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

// Git provides access to Azure DevOps Git REST APIs.
// This includes operations for repositories, commits, branches, and pull requests.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git
type Git struct {
	client Client
}

// RepoInfo retrieves detailed information about a Git repository.
// Returns the repository including its ID, default branch, URLs, and project reference.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get
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

// PRs returns a GitPRs client for pull request operations on the specified repository.
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

// GitPRs provides access to Azure DevOps Pull Request REST APIs.
// This includes operations for listing, creating, updating, and managing pull requests.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests
type GitPRs struct {
	client  Client
	baseUrl string
	org     string
}

// ByID retrieves a single pull request by its ID.
// Includes work item references and commits by default.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-request
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

// List retrieves pull requests matching the specified query criteria.
// Use SearchCriteria to filter by status, creator, reviewer, branches, etc.
//
// See:
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

// Create creates a new pull request.
// The GitPullRequest must include at minimum: SourceRefName, TargetRefName, and Title.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/create
func (g GitPRs) Create(
	ctx context.Context,
	pr models.GitPullRequest,
) (*models.GitPullRequest, error) {
	return httpPost[models.GitPullRequest](ctx, g.client, g.baseUrl, pr)
}

// Update modifies an existing pull request.
// Can update title, description, status, reviewers, and completion options.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/update
func (g GitPRs) Update(
	ctx context.Context,
	id int32,
	pr models.GitPullRequest,
) (*models.GitPullRequest, error) {
	prURL, _ := url.JoinPath(g.baseUrl, strconv.FormatInt(int64(id), 10))
	return httpPatch[models.GitPullRequest](ctx, g.client, prURL, pr)
}

// Reviewers return all reviewers for a pull request, including their votes.
//
// See:
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

// Vote sets the current user's vote on a pull request.
// Use models.PrVote constants for vote values:
//   - VoteApproved (10): Approve
//   - VoteApprovedWithSuggestions (5): Approve with suggestions
//   - VoteNone (0): Reset vote
//   - VoteWaitingForAuthor (-5): Waiting for author
//   - VoteRejected (-10): Reject
//
// See:
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

// Statuses return all statuses posted to a pull request.
// Statuses are typically posted by CI/CD systems to indicate build/test results.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-statuses/list
func (g GitPRs) Statuses(ctx context.Context, prID int32) ([]models.GitPullRequestStatus, error) {
	statusesURL, _ := url.JoinPath(
		g.baseUrl,
		strconv.FormatInt(int64(prID), 10),
		"statuses",
	)
	list, err := httpGet[List[models.GitPullRequestStatus]](ctx, g.client, statusesURL)
	if err != nil {
		return nil, err
	}
	return list.Value, nil
}

// List represents a paginated list response from Azure DevOps REST APIs.
// Most list endpoints return results in this wrapper format.
type List[T any] struct {
	// Value contains the list of items.
	Value []T `json:"value"`

	// Count is the number of items on the current page.
	Count int `json:"count"`
}
