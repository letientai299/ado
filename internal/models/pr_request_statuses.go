package models

// PullRequestStatus represents the status of the pull request.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type PullRequestStatus string

//goland:noinspection GoCommentStart,GoUnusedConst
const (
	// Status isn't set.
	// Default state.
	PullRequestStatusNotSet PullRequestStatus = "notSet"
	// Pull request is active.
	PullRequestStatusActive PullRequestStatus = "active"
	// Pull request is abandoned.
	PullRequestStatusAbandoned PullRequestStatus = "abandoned"
	// Pull request is completed.
	PullRequestStatusCompleted PullRequestStatus = "completed"
	// Used in pull request search criteria to include all statuses.
	PullRequestStatusAll PullRequestStatus = "all"
)
