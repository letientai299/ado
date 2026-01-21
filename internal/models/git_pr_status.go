package models

import "time"

// GitPullRequestStatus represents a status posted to a pull request.
// Statuses can be posted by external tools (like CI systems) to provide
// information about the state of builds, tests, or other checks.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-statuses/list#gitpullrequeststatus
type GitPullRequestStatus struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Id is the unique identifier of this status.
	Id int `json:"id,omitempty"`

	// State indicates the current state of the status.
	// Values: notSet, pending, succeeded, failed, error, notApplicable.
	State GitStatusState `json:"state,omitempty"`

	// Description is a human-readable description of the status.
	// Typically describes what the status represents or any error messages.
	Description string `json:"description,omitempty"`

	// Context uniquely identifies this status within the pull request.
	// Allows multiple statuses from different sources (e.g., different CI pipelines).
	Context GitStatusContext `json:"context"`

	// TargetUrl is a URL providing more details about the status.
	// Typically links to build results, test reports, or other external systems.
	TargetUrl string `json:"targetUrl,omitempty"`

	// IterationId is the PR iteration this status is associated with.
	// If not specified, the status applies to the entire pull request.
	IterationId int `json:"iterationId,omitempty"`

	// CreationDate is when this status was created.
	CreationDate *time.Time `json:"creationDate,omitempty"`

	// UpdatedDate is when this status was last updated.
	UpdatedDate *time.Time `json:"updatedDate,omitempty"`

	// CreatedBy is the identity that created this status.
	CreatedBy *IdentityRef `json:"createdBy,omitempty"`
}

// GitStatusContext uniquely identifies a status within a pull request or commit.
// The combination of Name and Genre forms a unique key for the status.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-statuses/create#gitstatuscontext
type GitStatusContext struct {
	// Name is the identifier of the status.
	// Cannot contain spaces. Examples: "build", "pr-validation", "jenkins/pr-merge".
	Name string `json:"name,omitempty"`

	// Genre is a category or grouping for the status.
	// Used for organizing and displaying statuses.
	// Common values: "continuous-integration", "testprtrigger", "prvalidation".
	Genre string `json:"genre,omitempty"`
}
