package models

// GitPullRequestMergeStrategy represents the strategy used to merge the pull
// request during completion.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitPullRequestMergeStrategy string

//goland:noinspection GoCommentStart,GoUnusedConst
const (
	// A two-parent, no-fast-forward merge.
	// The source branch is unchanged.
	// This is the default behavior.
	GitPullRequestMergeStrategyNoFastForward GitPullRequestMergeStrategy = "noFastForward"
	// Put all changes from the pull request into a single-parent commit.
	GitPullRequestMergeStrategySquash GitPullRequestMergeStrategy = "squash"
	// Rebase the source branch on top of the target branch HEAD commit, and
	// fast-forward the target branch.
	// The source branch is updated during the rebase operation.
	GitPullRequestMergeStrategyRebase GitPullRequestMergeStrategy = "rebase"
	// Rebase the source branch on top of the target branch HEAD commit, and create a
	// two-parent, no-fast-forward merge.
	// The source branch is updated during the rebase operation.
	GitPullRequestMergeStrategyRebaseMerge GitPullRequestMergeStrategy = "rebaseMerge"
)
