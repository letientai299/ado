package models

import "time"

// GitUserDate represents a user and a date.
type GitUserDate struct {
	Date  *time.Time `json:"date,omitempty"`
	Email string     `json:"email,omitempty"`
	Name  string     `json:"name,omitempty"`
}

// GitChangeCounts represents counts of changes by type.
type GitChangeCounts struct {
	Add    int `json:"add,omitempty"`
	Delete int `json:"delete,omitempty"`
	Edit   int `json:"edit,omitempty"`
}

// GitChange represents a change in a commit.
type GitChange struct {
	ChangeType string `json:"changeType,omitempty"`
	Item       any    `json:"item,omitempty"`
}

// GitStatus represents a status in a commit.
type GitStatus struct {
	Id          int    `json:"id,omitempty"`
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	Context     any    `json:"context,omitempty"`
}

// GitRepository represents a Git repository.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitRepository struct {
	// The class to represent a collection of REST reference links.
	Links         *ReferenceLinks `json:"_links,omitempty"`
	DefaultBranch string          `json:"defaultBranch,omitempty"`
	Id            string          `json:"id,omitempty"`
	// True if the repository is disabled.
	// False otherwise.
	IsDisabled bool `json:"isDisabled,omitempty"`
	// True if the repository was created as a fork.
	IsFork           bool           `json:"isFork,omitempty"`
	IsInMaintenance  bool           `json:"isInMaintenance,omitempty"`
	IsVariation      bool           `json:"isVariation,omitempty"`
	Name             string         `json:"name,omitempty"`
	ParentRepository *GitRepository `json:"parentRepository,omitempty"`
	Project          *TeamProject   `json:"project,omitempty"`
	RemoteUrl        string         `json:"remoteUrl,omitempty"`
	Size             int64          `json:"size,omitempty"`
	SshUrl           string         `json:"sshUrl,omitempty"`
	Url              string         `json:"url,omitempty"`
	ValidRemoteUrls  []string       `json:"validRemoteUrls,omitempty"`
	WebUrl           string         `json:"webUrl,omitempty"`
}

// GitCommitRef represents a reference to a Git commit.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitCommitRef struct {
	// A collection of REST reference links.
	Links *ReferenceLinks `json:"_links,omitempty"`
	// Author of the commit.
	Author *GitUserDate `json:"author,omitempty"`
	// Counts of changes by type.
	ChangeCounts *GitChangeCounts `json:"changeCounts,omitempty"`
	// An enumeration of the changes in a commit.
	Changes []GitChange `json:"changes,omitempty"`
	// Comment or message of the commit.
	Comment string `json:"comment,omitempty"`
	// Indicates if the comment is truncated from the full length.
	CommentTruncated bool `json:"commentTruncated,omitempty"`
	// Committer of the commit.
	Committer *GitUserDate `json:"committer,omitempty"`
	// SHA1 ID of the commit.
	CommitId string `json:"commitId,omitempty"`
	// The items in the commit.
	Statuses []GitStatus `json:"statuses,omitempty"`
	// Remote URL to the commit.
	Url string `json:"url,omitempty"`
	// Work item references for this commit.
	WorkItems []ResourceRef `json:"workItems,omitempty"`
}

// GitForkRef represents information about a fork source.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type GitForkRef struct {
	// The repository from which the fork was created.
	Repository *GitRepository `json:"repository,omitempty"`
	// The name of the branch in the source repository.
	RefName string `json:"refName,omitempty"`
	// The name of the repository from which the fork was created.
	Name string `json:"name,omitempty"`
	// The URL of the repository from which the fork was created.
	Url string `json:"url,omitempty"`
}
