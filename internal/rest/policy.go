package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
)

// Policy provides access to Azure DevOps Policy APIs
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy
type Policy struct {
	client Client
}

// Evaluations returns policy evaluations for a pull request
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/list
func (p Policy) Evaluations(
	ctx context.Context,
	repo config.Repository,
	projectID string,
	prID int32,
) ([]models.PolicyEvaluationRecord, error) {
	// Construct the artifact ID for the pull request
	// Format: vstfs:///CodeReview/CodeReviewId/{projectId}/{pullRequestId}
	artifactId := fmt.Sprintf("vstfs:///CodeReview/CodeReviewId/%s/%d", projectID, prID)

	// Policy evaluations API requires preview version
	apiURL := fmt.Sprintf("%s/%s/%s/_apis/policy/evaluations?artifactId=%s&api-version=7.1-preview.1",
		adoHost, repo.Org, repo.Project, url.QueryEscape(artifactId))

	// Make direct HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	list, err := call[List[models.PolicyEvaluationRecord]](p.client, req)
	if err != nil {
		return nil, err
	}

	return list.Value, nil
}
