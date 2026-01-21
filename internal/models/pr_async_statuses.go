package models

// PullRequestAsyncStatus represents the status of an asynchronous merge operation.
// Azure DevOps performs merges asynchronously and this status tracks the progress.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get#pullrequestasyncstatus
type PullRequestAsyncStatus string

const (
	// PullRequestAsyncStatusNotSet indicates the merge status has not been set.
	// This is the default/uninitialized state.
	PullRequestAsyncStatusNotSet PullRequestAsyncStatus = "notSet"

	// PullRequestAsyncStatusQueued indicates the merge operation is queued.
	// The merge is waiting to be processed.
	PullRequestAsyncStatusQueued PullRequestAsyncStatus = "queued"

	// PullRequestAsyncStatusConflicts indicates the merge has conflicts.
	// Manual conflict resolution is required before the PR can be completed.
	PullRequestAsyncStatusConflicts PullRequestAsyncStatus = "conflicts"

	// PullRequestAsyncStatusSucceeded indicates the merge completed successfully.
	// The source and target branches merged without conflicts.
	PullRequestAsyncStatusSucceeded PullRequestAsyncStatus = "succeeded"

	// PullRequestAsyncStatusRejectedByPolicy indicates the merge was rejected by policy.
	// One or more branch policies prevented the merge from completing.
	PullRequestAsyncStatusRejectedByPolicy PullRequestAsyncStatus = "rejectedByPolicy"

	// PullRequestAsyncStatusFailure indicates the merge operation failed.
	// An unexpected error occurred during the merge attempt.
	PullRequestAsyncStatusFailure PullRequestAsyncStatus = "failure"
)
