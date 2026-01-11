package models

// PullRequestAsyncStatus represents the current status of the pull request merge.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type PullRequestAsyncStatus string

//goland:noinspection GoCommentStart,GoUnusedConst
const (
	// Status is not set. Default state.
	PullRequestAsyncStatusNotSet PullRequestAsyncStatus = "notSet"
	// Pull request merge is queued.
	PullRequestAsyncStatusQueued PullRequestAsyncStatus = "queued"
	// Pull request merge failed due to conflicts.
	PullRequestAsyncStatusConflicts PullRequestAsyncStatus = "conflicts"
	// Pull request merge succeeded.
	PullRequestAsyncStatusSucceeded PullRequestAsyncStatus = "succeeded"
	// Pull request merge rejected by policy.
	PullRequestAsyncStatusRejectedByPolicy PullRequestAsyncStatus = "rejectedByPolicy"
	// Pull request merge failed.
	PullRequestAsyncStatusFailure PullRequestAsyncStatus = "failure"
)
