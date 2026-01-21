package models

// PrVote represents a reviewer's vote on a pull request.
// Votes indicate the reviewer's approval status and can block PR completion.
//
// Vote values use a numeric scale where positive values indicate approval
// and negative values indicate rejection or concerns.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-reviewers/list#identityrefwithvote
type PrVote int

const (
	// VoteApproved indicates the reviewer has approved the pull request.
	// Value: 10. This is a strong approval with no concerns.
	VoteApproved PrVote = 10

	// VoteApprovedWithSuggestions indicates approval with optional suggestions.
	// Value: 5. The reviewer approves but has non-blocking feedback.
	VoteApprovedWithSuggestions PrVote = 5

	// VoteNone indicates the reviewer has not voted.
	// Value: 0. This is the default state for new reviewers.
	VoteNone PrVote = 0

	// VoteWaitingForAuthor indicates the reviewer is waiting for changes.
	// Value: -5. The reviewer has concerns that should be addressed.
	VoteWaitingForAuthor PrVote = -5

	// VoteRejected indicates the reviewer has rejected the pull request.
	// Value: -10. This is a strong rejection that typically blocks completion.
	VoteRejected PrVote = -10
)
