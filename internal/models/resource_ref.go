package models

// ResourceRef represents a reference to a resource in Azure DevOps.
// This is a lightweight reference containing only the ID and URL of a resource,
// commonly used for linking to work items from commits or pull requests.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/commits/get#resourceref
type ResourceRef struct {
	// Id is the identifier of the referenced resource.
	// For work items, this is the work item ID as a string.
	Id string `json:"id,omitempty"`

	// Url is the REST API URL to retrieve the full resource.
	Url string `json:"url,omitempty"`
}
