package models

import "time"

// GitPullRequest represents all the data associated with a pull request.
// Pull requests let your team review code and give feedback on changes before
// merging into the main branch.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get#gitpullrequest
type GitPullRequest struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// ArtifactId uniquely identifies this pull request as an artifact.
	// Format: vstfs:///Git/PullRequestId/{projectId}/{repositoryId}/{pullRequestId}
	ArtifactId string `json:"artifactId,omitempty"`

	// AutoCompleteSetBy is the identity that enabled auto-complete for this PR.
	// If set, the PR will automatically complete when all policies pass.
	AutoCompleteSetBy *IdentityRef `json:"autoCompleteSetBy,omitempty"`

	// ClosedBy is the identity that closed the pull request.
	ClosedBy *IdentityRef `json:"closedBy,omitempty"`

	// ClosedDate is when the pull request was closed.
	// This is set when the PR is completed, abandoned, or merged externally.
	ClosedDate *time.Time `json:"closedDate,omitempty"`

	// CodeReviewId is the internal code review identifier.
	// Used internally by Azure DevOps.
	CodeReviewId int `json:"codeReviewId,omitempty"`

	// Commits contains the commits included in this pull request.
	// Only populated when includeCommits is true in the request.
	Commits []GitCommitRef `json:"commits,omitempty"`

	// CompletionOptions specifies how the PR will be merged when completed.
	CompletionOptions *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`

	// CompletionQueueTime is when the PR entered the completion queue.
	// Used internally by Azure DevOps.
	CompletionQueueTime *time.Time `json:"completionQueueTime,omitempty"`

	// CreatedBy is the identity that created this pull request.
	CreatedBy *IdentityRef `json:"createdBy,omitempty"`

	// CreationDate is when this pull request was created.
	CreationDate *time.Time `json:"creationDate,omitempty"`

	// Description is the detailed description of the pull request.
	// Supports markdown formatting.
	Description string `json:"description,omitempty"`

	// ForkSource contains information about the fork source if this PR
	// is from a forked repository.
	ForkSource *GitForkRef `json:"forkSource,omitempty"`

	// HasMultipleMergeBases indicates if there are multiple merge bases.
	// This can occur in complex branching scenarios and may require attention.
	HasMultipleMergeBases bool `json:"hasMultipleMergeBases"`

	// IsDraft indicates if this is a draft/work-in-progress pull request.
	// Draft PRs cannot be completed until marked as ready for review.
	IsDraft bool `json:"isDraft"`

	// Labels contains the labels/tags associated with this pull request.
	Labels []WebApiTagDefinition `json:"labels,omitempty"`

	// LastMergeCommit is the commit from the most recent successful merge.
	// Empty if the most recent merge is in progress or failed.
	LastMergeCommit *GitCommitRef `json:"lastMergeCommit,omitempty"`

	// LastMergeSourceCommit is the source branch HEAD at the time of the last merge.
	LastMergeSourceCommit *GitCommitRef `json:"lastMergeSourceCommit,omitempty"`

	// LastMergeTargetCommit is the target branch HEAD at the time of the last merge.
	LastMergeTargetCommit *GitCommitRef `json:"lastMergeTargetCommit,omitempty"`

	// MergeFailureMessage contains the error message if the merge failed.
	MergeFailureMessage string `json:"mergeFailureMessage,omitempty"`

	// MergeFailureType indicates the type of merge failure.
	MergeFailureType PullRequestMergeFailureType `json:"mergeFailureType,omitempty"`

	// MergeId is the internal identifier for the merge job.
	// Used internally by Azure DevOps.
	MergeId string `json:"mergeId,omitempty"`

	// MergeOptions specifies options for the merge operation.
	// These are separate from completion options and apply to each merge attempt.
	MergeOptions *GitPullRequestMergeOptions `json:"mergeOptions,omitempty"`

	// MergeStatus indicates the current status of the pull request merge.
	MergeStatus PullRequestAsyncStatus `json:"mergeStatus,omitempty"`

	// PullRequestId is the unique identifier of this pull request.
	PullRequestId int32 `json:"pullRequestId,omitempty"`

	// RemoteUrl is the internal remote URL.
	// Used internally by Azure DevOps.
	RemoteUrl string `json:"remoteUrl,omitempty"`

	// Repository is the repository containing the target branch.
	Repository *GitRepository `json:"repository,omitempty"`

	// Reviewers contains all reviewers and their votes on this PR.
	Reviewers []*IdentityRefWithVote `json:"reviewers,omitempty"`

	// SourceRefName is the full name of the source branch.
	// Format: refs/heads/{branch-name}
	SourceRefName string `json:"sourceRefName,omitempty"`

	// Status indicates the current state of the pull request.
	Status *PullRequestStatus `json:"status,omitempty"`

	// SupportsIterations indicates if this PR supports multiple iterations.
	// When true, each push to the source branch creates a new iteration,
	// allowing reviewers to track changes between iterations.
	SupportsIterations bool `json:"supportsIterations"`

	// TargetRefName is the full name of the target branch.
	// Format: refs/heads/{branch-name}
	TargetRefName string `json:"targetRefName,omitempty"`

	// Title is the title/subject of the pull request.
	Title string `json:"title,omitempty"`

	// Url is the REST API URL of this pull request.
	Url string `json:"url,omitempty"`

	// WorkItemRefs contains references to linked work items.
	WorkItemRefs []*ResourceRef `json:"workItemRefs,omitempty"`
}

// GitPullRequestCompletionOptions specifies how a pull request should be completed.
// These options control the merge behavior and post-completion actions.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/update#gitpullrequestcompletionoptions
type GitPullRequestCompletionOptions struct {
	// AutoCompleteIgnoreConfigIds lists policy configuration IDs that auto-complete
	// should not wait for. Only applies to optional (non-blocking) policies.
	AutoCompleteIgnoreConfigIds []int `json:"autoCompleteIgnoreConfigIds,omitempty"`

	// BypassPolicy bypasses all branch policies when completing the PR.
	// Requires elevated permissions and records the bypass reason.
	BypassPolicy bool `json:"bypassPolicy"`

	// BypassReason explains why policies were bypassed.
	// Required when BypassPolicy is true.
	BypassReason string `json:"bypassReason,omitempty"`

	// DeleteSourceBranch deletes the source branch after PR completion.
	DeleteSourceBranch bool `json:"deleteSourceBranch"`

	// MergeCommitMessage is the custom message for the merge commit.
	// If not specified, a default message is generated.
	MergeCommitMessage string `json:"mergeCommitMessage,omitempty"`

	// MergeStrategy specifies how to merge the PR.
	// Options: noFastForward (default), squash, rebase, rebaseMerge.
	// This supersedes the deprecated SquashMerge field.
	MergeStrategy GitPullRequestMergeStrategy `json:"mergeStrategy,omitempty"`

	// SquashMerge is deprecated. Use MergeStrategy instead.
	// When MergeStrategy is not set: false = no-fast-forward, true = squash.
	// When MergeStrategy is set, this field is ignored.
	SquashMerge bool `json:"squashMerge"`

	// TransitionWorkItems automatically transitions linked work items
	// to the next logical state (e.g., Active -> Resolved) upon completion.
	TransitionWorkItems bool `json:"transitionWorkItems"`

	// TriggeredByAutoComplete indicates if this completion was triggered
	// by auto-complete. Used internally by Azure DevOps.
	TriggeredByAutoComplete bool `json:"triggeredByAutoComplete"`
}

// GitPullRequestMergeOptions specifies options for the merge operation.
// These options apply each time a merge is attempted, not just at completion.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/update#gitpullrequestmergeoptions
type GitPullRequestMergeOptions struct {
	// ConflictAuthorshipCommits puts conflict resolutions in separate commits
	// to preserve authorship information for git blame.
	ConflictAuthorshipCommits bool `json:"conflictAuthorshipCommits"`

	// DetectRenameFalsePositives enables detection of false positive renames
	// during the merge operation.
	DetectRenameFalsePositives bool `json:"detectRenameFalsePositives"`

	// DisableRenames disables rename detection during the merge.
	// This can speed up merges for repositories with many files.
	DisableRenames bool `json:"disableRenames"`
}
