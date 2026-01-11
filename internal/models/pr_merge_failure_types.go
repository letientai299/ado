package models

// PullRequestMergeFailureType represents the type of failure (if any) of the pull request merge.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type PullRequestMergeFailureType string

//goland:noinspection GoCommentStart,GoUnusedConst
const (
	// Type is not set.
	// Default type.
	PullRequestMergeFailureTypeNone PullRequestMergeFailureType = "none"
	// Pull request merge failure type unknown.
	PullRequestMergeFailureTypeUnknown PullRequestMergeFailureType = "unknown"
	// Pull request merge failed due to case mismatch.
	PullRequestMergeFailureTypeCaseSensitive PullRequestMergeFailureType = "caseSensitive"
	// Pull request merge failed due to an object being too large.
	PullRequestMergeFailureTypeObjectTooLarge PullRequestMergeFailureType = "objectTooLarge"
)
