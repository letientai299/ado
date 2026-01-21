package models

// PullRequestTimeRangeType specifies which date field to use when filtering
// pull requests by time range using minTime and maxTime parameters.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests#pullrequesttimerangetype
type PullRequestTimeRangeType string

const (
	// PullRequestTimeRangeTypeCreated filters by the pull request creation date.
	// This is the default when no time range type is specified.
	// Use this to find PRs created within a specific time period.
	PullRequestTimeRangeTypeCreated PullRequestTimeRangeType = "created"

	// PullRequestTimeRangeTypeClosed filters by the pull request closed date.
	// Use this to find PRs that were completed or abandoned within a specific
	// time period. Only applies to non-active pull requests.
	PullRequestTimeRangeTypeClosed PullRequestTimeRangeType = "closed"
)
