package models

// WebApiTagDefinition represents a tag definition.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type WebApiTagDefinition struct {
	// Whether the tag is active.
	Active bool `json:"active,omitempty"`
	// The ID of the tag.
	Id string `json:"id,omitempty"`
	// The name of the tag.
	Name string `json:"name,omitempty"`
	// The URL of the tag.
	Url string `json:"url,omitempty"`
}
