package models

import "time"

// GitPullRequestStatus represents the status of a pull request at a particular iteration.
// Statuses can be posted by external tools (like CI systems) to provide information about the PR.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-statuses/list?view=azure-devops-rest-7.1
type GitPullRequestStatus struct {
	// ID of the status.
	Id int `json:"id,omitempty"`

	// Status state. Can be "pending", "succeeded", "failed", or "error".
	State string `json:"state,omitempty"`

	// Description of the status.
	Description string `json:"description,omitempty"`

	// Context of the status. Contains name and genre fields.
	Context GitStatusContext `json:"context,omitempty"`

	// Target URL for this status. This is typically a link to the build/test results.
	TargetUrl string `json:"targetUrl,omitempty"`

	// The iteration ID associated with the status.
	IterationId int `json:"iterationId,omitempty"`

	// Creation date of the status.
	CreationDate *time.Time `json:"creationDate,omitempty"`

	// Updated date of the status.
	UpdatedDate *time.Time `json:"updatedDate,omitempty"`

	// Identity that created the status.
	CreatedBy *IdentityRef `json:"createdBy,omitempty"`
}

// GitStatusContext uniquely identifies the status.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-statuses/list?view=azure-devops-rest-7.1
type GitStatusContext struct {
	// Name identifier of the status, cannot contain spaces.
	// For example, "continuous-integration/jenkins/pr-merge".
	Name string `json:"name,omitempty"`

	// Genre of the status. Typically used for grouping or display purposes.
	// Common values include "continuous-integration", "testprtrigger", "prvalidation".
	Genre string `json:"genre,omitempty"`
}
