package rest

import (
	"context"
	"fmt"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
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

// Evaluations retrieve policy evaluations for a pull request.
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
	prIDs ...int32,
) (map[int32][]models.PolicyEvaluationRecord, error) {
	if len(prIDs) == 0 {
		return nil, nil
	}
	ctx = WithAPIVersion(ctx, apiVersion7_1_preview)

	// Policy evaluations API requires preview version
	baseURL := fmt.Sprintf(
		"%s/%s/%s/_apis/policy/evaluations",
		adoHost,
		repo.Org,
		repo.Project,
	)

	// Fetch evaluations in parallel using goroutines
	type result struct {
		prID  int32
		evals []models.PolicyEvaluationRecord
		err   error
	}

	ch := make(chan result, len(prIDs))
	for _, id := range prIDs {
		go func(prID int32) {
			artifactId := fmt.Sprintf("vstfs:///CodeReview/CodeReviewId/%s/%d", projectID, prID)
			list, err := httpGet[List[models.PolicyEvaluationRecord]](
				ctx,
				p.client,
				baseURL,
				_shared.KV[string]{Key: "artifactId", Value: artifactId},
			)
			if err != nil {
				ch <- result{prID: prID, err: err}
				return
			}
			ch <- result{prID: prID, evals: list.Value}
		}(id)
	}

	finalResult := make(map[int32][]models.PolicyEvaluationRecord, len(prIDs))
	for range prIDs {
		res := <-ch
		if res.err != nil {
			return nil, res.err
		}
		finalResult[res.prID] = res.evals
	}

	return finalResult, nil
}

// Requeue re-evaluates a policy for a pull request.
// Use this to trigger a fresh evaluation of policies, for example, after
// fixing a build validation failure.
//
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/requeue-policy-evaluation
func (p Policy) Requeue(
	ctx context.Context,
	repo config.Repository,
	evaluationID string,
) (*models.PolicyEvaluationRecord, error) {
	apiURL := fmt.Sprintf(
		"%s/%s/%s/_apis/policy/evaluations/%s",
		adoHost,
		repo.Org,
		repo.Project,
		evaluationID,
	)
	return httpPatch[models.PolicyEvaluationRecord](ctx, p.client, apiURL, nil)
}
