package models

// GitPullRequestMergeStrategy specifies how the pull request will be merged
// when it is completed. Different strategies have different effects on the
// commit history.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/update#gitpullrequestmergestrategy
type GitPullRequestMergeStrategy string

const (
	// GitPullRequestMergeStrategyNoFastForward creates a merge commit.
	// This is the default strategy. Creates a two-parent merge commit even if
	// a fast-forward would be possible. The source branch remains unchanged.
	// Preserves the complete branch history.
	GitPullRequestMergeStrategyNoFastForward GitPullRequestMergeStrategy = "noFastForward"

	// GitPullRequestMergeStrategySquash squashes all commits into one.
	// Combines all commits from the pull request into a single commit on the
	// target branch. Useful for keeping the target branch history clean.
	GitPullRequestMergeStrategySquash GitPullRequestMergeStrategy = "squash"

	// GitPullRequestMergeStrategyRebase rebases and fast-forwards.
	// Rebases the source branch commits on top of the target branch HEAD,
	// then fast-forwards the target branch. The source branch is updated
	// during the rebase. Results in a linear history.
	GitPullRequestMergeStrategyRebase GitPullRequestMergeStrategy = "rebase"

	// GitPullRequestMergeStrategyRebaseMerge rebases then creates a merge commit.
	// Rebases the source branch on top of the target branch HEAD, then creates
	// a two-parent merge commit. The source branch is updated during rebase.
	// Combines rebase benefits with merge commit visibility.
	GitPullRequestMergeStrategyRebaseMerge GitPullRequestMergeStrategy = "rebaseMerge"
)
