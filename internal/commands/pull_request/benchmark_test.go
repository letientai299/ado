package pull_request

import (
	"testing"
	"time"

	"github.com/letientai299/ado/internal/models"
)

// BenchmarkConverterWithStatuses benchmarks the PR conversion function.
// This is used in both `pr list`, `pr view`, and `pr update`.
func BenchmarkConverterWithStatuses(b *testing.B) {
	baseURL := "https://dev.azure.com/org/project/_git/repo/pullrequest"
	orgName := "testorg"
	repo := &models.GitRepository{
		Id:   "repo-id",
		Name: "repo",
		Project: &models.TeamProject{
			Id:   "project-id",
			Name: "project",
		},
	}

	tests := []struct {
		name        string
		pr          models.GitPullRequest
		evaluations map[int32][]models.PolicyEvaluationRecord
	}{
		{
			name:        "simple_pr_no_evals",
			pr:          makePR(1, 2),
			evaluations: nil,
		},
		{
			name:        "pr_with_5_reviewers",
			pr:          makePR(5, 0),
			evaluations: nil,
		},
		{
			name:        "pr_with_evaluations",
			pr:          makePR(3, 2),
			evaluations: makeEvaluations(1, 5),
		},
		{
			name:        "pr_with_many_evaluations",
			pr:          makePR(5, 3),
			evaluations: makeEvaluations(1, 20),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			converter := converter(baseURL, orgName, repo, tt.evaluations)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = converter(tt.pr)
			}
		})
	}
}

// BenchmarkResolvePolicyChecks benchmarks policy evaluation processing.
func BenchmarkResolvePolicyChecks(b *testing.B) {
	tests := []struct {
		name        string
		numEvals    int
		hasConflict bool
	}{
		{"no_evals", 0, false},
		{"5_evals", 5, false},
		{"10_evals_with_conflict", 10, true},
		{"20_evals", 20, false},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			pr := makePR(2, 1)
			if tt.hasConflict {
				pr.MergeStatus = models.PullRequestAsyncStatusConflicts
			} else {
				pr.MergeStatus = models.PullRequestAsyncStatusSucceeded
			}
			evals := makeEvalsList(tt.numEvals)
			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_ = resolvePolicyChecks(&pr, evals)
			}
		})
	}
}

// BenchmarkPolicyChecksSummary benchmarks the policy summary generation.
func BenchmarkPolicyChecksSummary(b *testing.B) {
	tests := []struct {
		name   string
		checks PolicyChecks
	}{
		{"empty", nil},
		{"all_passed", makePolicyChecks(5, "approved")},
		{"some_failed", makeMixedPolicyChecks(5)},
		{"many_checks", makeMixedPolicyChecks(20)},
	}

	for _, tt := range tests {
		b.Run(tt.name+"_text", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tt.checks.SummaryText()
			}
		})
		b.Run(tt.name+"_icon", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = tt.checks.SummaryIcon()
			}
		})
	}
}

func makePR(numReviewers, numApprovers int) models.GitPullRequest {
	now := time.Now()
	reviewers := make([]*models.IdentityRefWithVote, numReviewers)
	for i := range numReviewers {
		vote := 0
		if i < numApprovers {
			vote = 10 // approved
		}
		reviewers[i] = &models.IdentityRefWithVote{
			IdentityRef: models.IdentityRef{
				Id:          "reviewer-" + string(rune('0'+i)),
				DisplayName: "Reviewer " + string(rune('A'+i)),
				UniqueName:  "reviewer" + string(rune('0'+i)) + "@example.com",
			},
			Vote: vote,
		}
	}

	return models.GitPullRequest{
		PullRequestId: 123,
		Title:         "Test PR Title",
		Description:   "This is a test PR description with some content.",
		IsDraft:       false,
		CreatedBy: &models.IdentityRef{
			Id:          "author-id",
			DisplayName: "Author Name",
			UniqueName:  "author@example.com",
		},
		CreationDate:  &now,
		SourceRefName: "refs/heads/feature/test",
		TargetRefName: "refs/heads/main",
		Reviewers:     reviewers,
		MergeStatus:   models.PullRequestAsyncStatusSucceeded,
	}
}

func makeEvaluations(prID int32, numEvals int) map[int32][]models.PolicyEvaluationRecord {
	return map[int32][]models.PolicyEvaluationRecord{
		prID: makeEvalsList(numEvals),
	}
}

func makeEvalsList(numEvals int) []models.PolicyEvaluationRecord {
	evals := make([]models.PolicyEvaluationRecord, numEvals)
	statuses := []models.PolicyEvaluationStatus{
		models.PolicyEvaluationStatusApproved,
		models.PolicyEvaluationStatusRejected,
		models.PolicyEvaluationStatusRunning,
		models.PolicyEvaluationStatusQueued,
	}
	for i := range numEvals {
		evals[i] = models.PolicyEvaluationRecord{
			EvaluationId: "eval-" + string(rune('0'+i)),
			Status:       statuses[i%len(statuses)],
			Configuration: models.PolicyConfiguration{
				IsBlocking: i%2 == 0,
				Type: models.PolicyTypeRef{
					Id:          "policy-type-" + string(rune('0'+i)),
					DisplayName: "Policy " + string(rune('A'+i)),
				},
				Settings: map[string]any{
					"displayName": "Policy Display " + string(rune('A'+i)),
				},
			},
			Context: map[string]any{
				"buildId": 1000 + i,
			},
		}
	}
	return evals
}

func makePolicyChecks(n int, status string) PolicyChecks {
	checks := make(PolicyChecks, n)
	for i := range n {
		checks[i] = PolicyCheck{
			Name:       "Check " + string(rune('A'+i)),
			Status:     status,
			IsRequired: true,
			Icon:       "✓",
		}
	}
	return checks
}

func makeMixedPolicyChecks(n int) PolicyChecks {
	checks := make(PolicyChecks, n)
	statuses := []string{"approved", "rejected", "running", "queued"}
	for i := range n {
		checks[i] = PolicyCheck{
			Name:       "Check " + string(rune('A'+i%26)),
			Status:     statuses[i%len(statuses)],
			IsRequired: i%3 != 0, // 2/3 are required
			Icon:       "?",
		}
	}
	return checks
}
