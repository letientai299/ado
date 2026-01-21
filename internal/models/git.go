package models

import "time"

// VersionControlChangeType represents the type of change made to an item in a commit.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#versioncontrolchangetype
type VersionControlChangeType string

const (
	// VersionControlChangeTypeNone indicates no change.
	VersionControlChangeTypeNone VersionControlChangeType = "none"
	// VersionControlChangeTypeAdd indicates the item was added.
	VersionControlChangeTypeAdd VersionControlChangeType = "add"
	// VersionControlChangeTypeEdit indicates the item was edited/modified.
	VersionControlChangeTypeEdit VersionControlChangeType = "edit"
	// VersionControlChangeTypeEncoding indicates the item's encoding changed.
	VersionControlChangeTypeEncoding VersionControlChangeType = "encoding"
	// VersionControlChangeTypeRename indicates the item was renamed.
	VersionControlChangeTypeRename VersionControlChangeType = "rename"
	// VersionControlChangeTypeDelete indicates the item was deleted.
	VersionControlChangeTypeDelete VersionControlChangeType = "delete"
	// VersionControlChangeTypeUndelete indicates the item was undeleted.
	VersionControlChangeTypeUndelete VersionControlChangeType = "undelete"
	// VersionControlChangeTypeBranch indicates a branch operation.
	VersionControlChangeTypeBranch VersionControlChangeType = "branch"
	// VersionControlChangeTypeMerge indicates a merge operation.
	VersionControlChangeTypeMerge VersionControlChangeType = "merge"
	// VersionControlChangeTypeLock indicates the item was locked.
	VersionControlChangeTypeLock VersionControlChangeType = "lock"
	// VersionControlChangeTypeRollback indicates a rollback operation.
	VersionControlChangeTypeRollback VersionControlChangeType = "rollback"
	// VersionControlChangeTypeSourceRename indicates the source was renamed.
	VersionControlChangeTypeSourceRename VersionControlChangeType = "sourceRename"
	// VersionControlChangeTypeTargetRename indicates the target was renamed.
	VersionControlChangeTypeTargetRename VersionControlChangeType = "targetRename"
	// VersionControlChangeTypeProperty indicates a property change.
	VersionControlChangeTypeProperty VersionControlChangeType = "property"
	// VersionControlChangeTypeAll indicates all change types.
	VersionControlChangeTypeAll VersionControlChangeType = "all"
)

// GitStatusState represents the state of a Git status.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/statuses/create#gitstatusstate
type GitStatusState string

const (
	// GitStatusStateNotSet indicates the status state is not set.
	GitStatusStateNotSet GitStatusState = "notSet"
	// GitStatusStatePending indicates the status is pending.
	GitStatusStatePending GitStatusState = "pending"
	// GitStatusStateSucceeded indicates the status succeeded.
	GitStatusStateSucceeded GitStatusState = "succeeded"
	// GitStatusStateFailed indicates the status failed.
	GitStatusStateFailed GitStatusState = "failed"
	// GitStatusStateError indicates the status encountered an error.
	GitStatusStateError GitStatusState = "error"
	// GitStatusStateNotApplicable indicates the status is not applicable.
	GitStatusStateNotApplicable GitStatusState = "notApplicable"
)

// GitUserDate represents a Git user with associated timestamp information.
// This is commonly used to identify the author or committer of a Git commit.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#gituserdate
type GitUserDate struct {
	// Date is the timestamp when the action occurred.
	Date *time.Time `json:"date,omitempty"`

	// Email is the email address of the user.
	Email string `json:"email,omitempty"`

	// Name is the display name of the user.
	Name string `json:"name,omitempty"`

	// ImageUrl is the URL of the user's avatar image.
	ImageUrl string `json:"imageUrl,omitempty"`
}

// GitChangeCounts represents counts of changes by type in a commit.
// These counts summarize the number of additions, deletions, and edits.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#changecounts
type GitChangeCounts struct {
	// Add is the number of files added.
	Add int `json:"add,omitempty"`

	// Delete is the number of files deleted.
	Delete int `json:"delete,omitempty"`

	// Edit is the number of files edited/modified.
	Edit int `json:"edit,omitempty"`
}

// GitChange represents a change to a file or folder in a Git commit.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#gitchange
type GitChange struct {
	// ChangeId is the unique identifier for this change.
	ChangeId int `json:"changeId,omitempty"`

	// ChangeType indicates what kind of change was made (add, edit, delete, etc.).
	ChangeType VersionControlChangeType `json:"changeType,omitempty"`

	// Item contains information about the item that was changed.
	// This is typically a GitItem object.
	Item any `json:"item,omitempty"`

	// NewContent contains the new content of the item (for adds and edits).
	NewContent *ItemContent `json:"newContent,omitempty"`

	// OriginalPath is the original path of the item before rename/move.
	OriginalPath string `json:"originalPath,omitempty"`

	// SourceServerItem is the server path of the source item.
	SourceServerItem string `json:"sourceServerItem,omitempty"`

	// Url is the REST API URL for this change.
	Url string `json:"url,omitempty"`
}

// ItemContent represents the content of an item.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#itemcontent
type ItemContent struct {
	// Content is the actual content of the item.
	Content string `json:"content,omitempty"`

	// ContentType indicates the type of content (rawText, base64Encoded, etc.).
	ContentType string `json:"contentType,omitempty"`
}

// GitStatus represents a status associated with a Git ref (commit, branch, etc.).
// Statuses are typically posted by external systems like CI/CD pipelines to indicate
// the state of builds or checks.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/statuses/list#gitstatus
type GitStatus struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Id is the unique identifier of the status.
	Id int `json:"id,omitempty"`

	// State indicates the current state of the status
	// (pending, succeeded, failed, error, notApplicable, notSet).
	State GitStatusState `json:"state,omitempty"`

	// Description is a human-readable description of the status.
	Description string `json:"description,omitempty"`

	// Context uniquely identifies the status within a ref.
	// Contains name and genre fields.
	Context *GitStatusContext `json:"context,omitempty"`

	// TargetUrl is a URL that provides more details about the status.
	// Typically links to build results or test reports.
	TargetUrl string `json:"targetUrl,omitempty"`

	// CreationDate is when the status was created.
	CreationDate *time.Time `json:"creationDate,omitempty"`

	// UpdatedDate is when the status was last updated.
	UpdatedDate *time.Time `json:"updatedDate,omitempty"`

	// CreatedBy is the identity that created this status.
	CreatedBy *IdentityRef `json:"createdBy,omitempty"`
}

// GitRepository represents a Git repository in Azure DevOps.
// Repositories contain source code, branches, commits, and pull requests.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get#gitrepository
type GitRepository struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// DefaultBranch is the name of the default branch (e.g., "refs/heads/main").
	DefaultBranch string `json:"defaultBranch,omitempty"`

	// Id is the unique identifier (GUID) of the repository.
	Id string `json:"id,omitempty"`

	// IsDisabled indicates whether the repository is disabled.
	// Disabled repositories cannot be accessed.
	IsDisabled bool `json:"isDisabled,omitempty"`

	// IsFork indicates whether this repository was created as a fork
	// of another repository.
	IsFork bool `json:"isFork,omitempty"`

	// IsInMaintenance indicates whether the repository is in maintenance mode.
	// Repositories in maintenance mode have limited functionality.
	IsInMaintenance bool `json:"isInMaintenance,omitempty"`

	// IsVariation indicates whether this is a variation of another repository.
	IsVariation bool `json:"isVariation,omitempty"`

	// Name is the name of the repository.
	Name string `json:"name,omitempty"`

	// ParentRepository is the parent repository if this is a fork.
	ParentRepository *GitRepository `json:"parentRepository,omitempty"`

	// Project is the team project that contains this repository.
	Project *TeamProject `json:"project,omitempty"`

	// RemoteUrl is the HTTPS URL for cloning this repository.
	RemoteUrl string `json:"remoteUrl,omitempty"`

	// Size is the size of the repository in bytes.
	Size int64 `json:"size,omitempty"`

	// SshUrl is the SSH URL for cloning this repository.
	SshUrl string `json:"sshUrl,omitempty"`

	// Url is the REST API URL of this repository.
	Url string `json:"url,omitempty"`

	// ValidRemoteUrls contains valid remote URLs for this repository.
	ValidRemoteUrls []string `json:"validRemoteUrls,omitempty"`

	// WebUrl is the URL to view this repository in a web browser.
	WebUrl string `json:"webUrl,omitempty"`
}

// GitCommitRef represents a reference to a Git commit with optional details.
// This includes commit metadata like author, committer, message, and change statistics.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#gitcommitref
type GitCommitRef struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Author is the person who authored the commit.
	// The author is the person who originally wrote the code.
	Author *GitUserDate `json:"author,omitempty"`

	// ChangeCounts contains counts of files added, deleted, and edited.
	ChangeCounts *GitChangeCounts `json:"changeCounts,omitempty"`

	// Changes is the list of changes included in the commit.
	Changes []GitChange `json:"changes,omitempty"`

	// Comment is the commit message.
	Comment string `json:"comment,omitempty"`

	// CommentTruncated indicates if the comment was truncated.
	// Long commit messages may be truncated for performance.
	CommentTruncated bool `json:"commentTruncated,omitempty"`

	// Committer is the person who committed the code.
	// The committer may differ from the author (e.g., when applying patches).
	Committer *GitUserDate `json:"committer,omitempty"`

	// CommitId is the full SHA-1 hash of the commit.
	CommitId string `json:"commitId,omitempty"`

	// Parents contains the parent commit IDs.
	// Merge commits have multiple parents.
	Parents []string `json:"parents,omitempty"`

	// Push contains information about the push that introduced this commit.
	Push *GitPushRef `json:"push,omitempty"`

	// RemoteUrl is the URL to view this commit in a web browser.
	RemoteUrl string `json:"remoteUrl,omitempty"`

	// Statuses contains Git statuses associated with this commit.
	Statuses []GitStatus `json:"statuses,omitempty"`

	// Url is the REST API URL of this commit.
	Url string `json:"url,omitempty"`

	// WorkItems contains references to work items linked to this commit.
	WorkItems []ResourceRef `json:"workItems,omitempty"`
}

// GitPushRef represents a reference to a Git push operation.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pushes/get#gitpushref
type GitPushRef struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Date is when the push occurred.
	Date *time.Time `json:"date,omitempty"`

	// PushId is the unique identifier of the push.
	PushId int `json:"pushId,omitempty"`

	// PushedBy is the identity that performed the push.
	PushedBy *IdentityRef `json:"pushedBy,omitempty"`

	// Url is the REST API URL of this push.
	Url string `json:"url,omitempty"`
}

// GitForkRef represents information about a fork source.
// This is used when a pull request originates from a forked repository.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/forks/list#gitforkref
type GitForkRef struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Creator is the identity that created the fork.
	Creator *IdentityRef `json:"creator,omitempty"`

	// IsLocked indicates whether the ref is locked.
	IsLocked bool `json:"isLocked,omitempty"`

	// IsLockedBy is the identity that locked the ref.
	IsLockedBy *IdentityRef `json:"isLockedBy,omitempty"`

	// Name is the name of the ref (e.g., "refs/heads/feature-branch").
	Name string `json:"name,omitempty"`

	// ObjectId is the commit ID that the ref points to.
	ObjectId string `json:"objectId,omitempty"`

	// PeeledObjectId is the ID of the object the ref points to (for tags).
	PeeledObjectId string `json:"peeledObjectId,omitempty"`

	// Repository is the forked repository.
	Repository *GitRepository `json:"repository,omitempty"`

	// Statuses contains Git statuses associated with this ref.
	Statuses []GitStatus `json:"statuses,omitempty"`

	// Url is the REST API URL of this ref.
	Url string `json:"url,omitempty"`
}
