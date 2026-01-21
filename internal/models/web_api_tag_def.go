package models

// WebApiTagDefinition represents a tag/label definition in Azure DevOps.
// Tags can be applied to pull requests, work items, and other resources
// for categorization and filtering.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-labels/list#webapitagdefinition
type WebApiTagDefinition struct {
	// Active indicates whether the tag is currently active.
	// Inactive tags may not appear in UI but still exist.
	Active bool `json:"active,omitempty"`

	// Id is the unique identifier (GUID) of the tag.
	Id string `json:"id,omitempty"`

	// Name is the display name of the tag.
	Name string `json:"name,omitempty"`

	// Url is the REST API URL of the tag definition.
	Url string `json:"url,omitempty"`
}
