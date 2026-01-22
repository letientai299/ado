package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
)

// Policy provides access to Azure DevOps Policy APIs.
// This client wraps the Policy REST API for working with branch policies
// and policy evaluations.
//
// Branch policies enforce code quality standards on pull requests, such as
// requiring builds to pass, minimum reviewer counts, or work item linking.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/policy
type Policy struct {
	client Client
}

// Evaluations retrieves policy evaluations for a pull request.
// Policy evaluations show the status of each branch policy for the given PR,
// including build validation policies, required reviewers, and other checks.
//
// The projectID should be the GUID of the project (available from repo.Project.Id).
// The prID is the numeric pull request ID.
//
// Returns a list of [models.PolicyEvaluationRecord], one for each policy
// configured for the target branch. Use [models.PolicyTypeBuildValidation]
// to identify build validation policies.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/list
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
	apiURL := fmt.Sprintf(
		"%s/%s/%s/_apis/policy/evaluations?artifactId=%s&api-version=7.1-preview.1",
		adoHost,
		repo.Org,
		repo.Project,
		url.QueryEscape(artifactId),
	)

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

// Requeue re-evaluates a policy for a pull request.
// Use this to trigger a fresh evaluation of policies, for example after
// fixing a build validation failure.
//
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/requeue-policy-evaluation
func (p Policy) Requeue(
	ctx context.Context,
	repo config.Repository,
	projectID string,
	evaluationID string,
) (*models.PolicyEvaluationRecord, error) {
	apiURL := fmt.Sprintf(
		"%s/%s/%s/_apis/policy/evaluations/%s?api-version=7.1-preview.1",
		adoHost,
		repo.Org,
		repo.Project,
		evaluationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, apiURL, nil)
	if err != nil {
		return nil, err
	}

	return call[models.PolicyEvaluationRecord](p.client, req)
}
