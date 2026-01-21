package models

// PullRequestStatus represents the lifecycle state of a pull request.
// Pull requests transition through these states from creation to completion.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get#pullrequeststatus
type PullRequestStatus string

const (
	// PullRequestStatusNotSet indicates the status has not been set.
	// This is the default/uninitialized state.
	PullRequestStatusNotSet PullRequestStatus = "notSet"

	// PullRequestStatusActive indicates the pull request is open and active.
	// Active PRs are awaiting review, approval, or completion.
	PullRequestStatusActive PullRequestStatus = "active"

	// PullRequestStatusAbandoned indicates the pull request was abandoned.
	// Abandoned PRs are closed without merging and cannot be completed.
	PullRequestStatusAbandoned PullRequestStatus = "abandoned"

	// PullRequestStatusCompleted indicates the pull request was completed.
	// Completed PRs have been successfully merged into the target branch.
	PullRequestStatusCompleted PullRequestStatus = "completed"

	// PullRequestStatusAll is used in search criteria to match all statuses.
	// This is not a valid PR state, only used for filtering.
	PullRequestStatusAll PullRequestStatus = "all"
)
