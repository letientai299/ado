package models

// PrVote represents a reviewer's vote on a pull request.
// https://learn.microsoft.com/en-us/javascript/api/azure-devops-extension-api/identityrefwithvote#azure-devops-extension-api-identityrefwithvote-vote
type PrVote int

const (
	VoteApproved                PrVote = 10
	VoteApprovedWithSuggestions PrVote = 5
	VoteNone                    PrVote = 0
	VoteWaitingForAuthor        PrVote = -5
	VoteRejected                PrVote = -10
)
