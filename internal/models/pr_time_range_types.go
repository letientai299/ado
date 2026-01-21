package models

// PullRequestTimeRangeType is the type of time range which should be used for minTime and maxTime.
// Defaults to PullRequestTimeRangeTypeCreated if unset.
//
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests?view=azure-devops-rest-7.1&tabs=HTTP#pullrequesttimerangetype
//
//goland:noinspection GoCommentStart,GoUnusedConst
type PullRequestTimeRangeType string

//goland:noinspection GoCommentStart,GoUnusedConst
const (
	PullRequestTimeRangeTypeCreated PullRequestTimeRangeType = "created"
	PullRequestTimeRangeTypeClosed  PullRequestTimeRangeType = "closed"
)
