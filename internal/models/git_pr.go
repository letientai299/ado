package models

import "time"

// GitPullRequest represents all the data associated with a pull request.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitPullRequest struct {
	// Links to other related objects.
	Links *ReferenceLinks `json:"_links,omitempty"`
	// A string which uniquely identifies this pull request.
	// To generate an artifact ID for a pull request, use this template:
	// vstfs:///Git/PullRequestId/{projectId}/{repositoryId}/{pullRequestId}
	ArtifactId string `json:"artifactId,omitempty"`
	// If set, auto-complete is enabled for this pull request, and this is the
	// identity that enabled it.
	AutoCompleteSetBy *IdentityRef `json:"autoCompleteSetBy,omitempty"`
	// The user who closed the pull request.
	ClosedBy *IdentityRef `json:"closedBy,omitempty"`
	// The date when the pull request was closed (completed, abandoned, or merged
	// externally).
	ClosedDate *time.Time `json:"closedDate,omitempty"`
	// The code review ID of the pull request.
	// Used internally.
	CodeReviewId int `json:"codeReviewId,omitempty"`
	// The commits contained in the pull request.
	Commits []GitCommitRef `json:"commits,omitempty"`
	// Options which affect how the pull request will be merged when it is
	// completed.
	CompletionOptions *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`
	// The most recent date at which the pull request entered the queue to be
	// completed.
	// Used internally.
	CompletionQueueTime *time.Time `json:"completionQueueTime,omitempty"`
	// The identity of the user who created the pull request.
	CreatedBy *IdentityRef `json:"createdBy,omitempty"`
	// The date when the pull request was created.
	CreationDate *time.Time `json:"creationDate,omitempty"`
	// The description of the pull request.
	Description string `json:"description,omitempty"`
	// If this is a PR from a fork, this will contain information about its source.
	ForkSource *GitForkRef `json:"forkSource,omitempty"`
	// Multiple merge bases warning
	HasMultipleMergeBases bool `json:"hasMultipleMergeBases,omitempty"`
	// Draft / WIP pull request.
	IsDraft bool `json:"isDraft,omitempty"`
	// The labels associated with the pull request.
	Labels []WebApiTagDefinition `json:"labels,omitempty"`
	// The commit of the most recent pull request merge.
	// If empty, the most recent merge is in progress or was unsuccessful.
	LastMergeCommit *GitCommitRef `json:"lastMergeCommit,omitempty"`
	// The commit at the head of the source branch at the time of the last pull
	// request merge.
	LastMergeSourceCommit *GitCommitRef `json:"lastMergeSourceCommit,omitempty"`
	// The commit at the head of the target branch at the time of the last pull
	// request merge.
	LastMergeTargetCommit *GitCommitRef `json:"lastMergeTargetCommit,omitempty"`
	// If set, the pull request merge failed for this reason.
	MergeFailureMessage string `json:"mergeFailureMessage,omitempty"`
	// The type of failure (if any) of the pull request merge.
	MergeFailureType PullRequestMergeFailureType `json:"mergeFailureType,omitempty"`
	// The ID of the job used to run the pull request merge.
	// Used internally.
	MergeId string `json:"mergeId,omitempty"`
	// Options used when the pull request merge runs.
	// These are separate from completion options since completion happens only once
	// and a new merge will run every time the source branch of the pull request
	// changes.
	MergeOptions *GitPullRequestMergeOptions `json:"mergeOptions,omitempty"`
	// The current status of the pull request merge.
	MergeStatus PullRequestAsyncStatus `json:"mergeStatus,omitempty"`
	// The ID of the pull request.
	PullRequestId int `json:"pullRequestId,omitempty"`
	// Used internally.
	RemoteUrl string `json:"remoteUrl,omitempty"`
	// The repository containing the target branch of the pull request.
	Repository *GitRepository `json:"repository,omitempty"`
	// A list of reviewers on the pull request along with the state of their votes.
	Reviewers []IdentityRefWithVote `json:"reviewers,omitempty"`
	// The name of the source branch of the pull request.
	SourceRefName string `json:"sourceRefName,omitempty"`
	// The status of the pull request.
	Status PullRequestStatus `json:"status,omitempty"`
	// If true, this pull request supports multiple iterations.
	// Iteration support means individual pushes to the source branch of the pull
	// request can be reviewed and comments left in one iteration will be tracked
	// across future iterations.
	SupportsIterations bool `json:"supportsIterations,omitempty"`
	// The name of the target branch of the pull request.
	TargetRefName string `json:"targetRefName,omitempty"`
	// The title of the pull request.
	Title string `json:"title,omitempty"`
	// Used internally.
	Url string `json:"url,omitempty"`
	// Any work item references associated with this pull request.
	WorkItemRefs []ResourceRef `json:"workItemRefs,omitempty"`
}

// GitPullRequestCompletionOptions represents preferences about how the pull request should be
// completed.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitPullRequestCompletionOptions struct {
	// List of any policy configuration ID's which auto-complete should not wait for.
	// Only applies to optional policies (isBlocking == false).
	// Auto-complete always waits for required policies (isBlocking == true).
	AutoCompleteIgnoreConfigIds []int `json:"autoCompleteIgnoreConfigIds,omitempty"`
	// If true, policies will be explicitly bypassed while the pull request is
	// completed.
	BypassPolicy bool `json:"bypassPolicy,omitempty"`
	// If policies are bypassed, this reason is stored as to why bypass was used.
	BypassReason string `json:"bypassReason,omitempty"`
	// If true, the source branch of the pull request will be deleted after
	// completion.
	DeleteSourceBranch bool `json:"deleteSourceBranch,omitempty"`
	// If set, this will be used as the commit message of the merge commit.
	MergeCommitMessage string `json:"mergeCommitMessage,omitempty"`
	// Specify the strategy used to merge the pull request during completion.
	// If MergeStrategy is not set to any value, a no-FF merge will be created if
	// SquashMerge == false.
	// If MergeStrategy is not set to any value, the pull request commits will be
	// squashed if SquashMerge == true.
	// The SquashMerge property is deprecated.
	// It is recommended that you explicitly set MergeStrategy in all cases.
	// If an explicit value is provided for MergeStrategy, the SquashMerge property
	// will be ignored.
	MergeStrategy GitPullRequestMergeStrategy `json:"mergeStrategy,omitempty"`
	// SquashMerge is deprecated.
	// You should explicitly set the value of MergeStrategy.
	// If MergeStrategy is set to any value, the SquashMerge value will be ignored.
	// If MergeStrategy is not set, the merge strategy will be no-fast-forward if
	// this flag is false, or squash if true.
	SquashMerge bool `json:"squashMerge,omitempty"`
	// If true, we will attempt to transition any work items linked to the pull
	// request into the next logical state (i.e.
	// Active -> Resolved)
	TransitionWorkItems bool `json:"transitionWorkItems,omitempty"`
	// If true, the current completion attempt was triggered via auto-complete.
	// Used internally.
	TriggeredByAutoComplete bool `json:"triggeredByAutoComplete,omitempty"`
}

// GitPullRequestMergeOptions represents the options which are used when a pull request merge is
// created.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitPullRequestMergeOptions struct {
	// If true, conflict resolutions applied during the merge will be put in
	// separate commits to preserve authorship info for git blame, etc.
	ConflictAuthorshipCommits  bool `json:"conflictAuthorshipCommits,omitempty"`
	DetectRenameFalsePositives bool `json:"detectRenameFalsePositives,omitempty"`
	// If true, rename detection will not be performed during the merge.
	DisableRenames bool `json:"disableRenames,omitempty"`
}
