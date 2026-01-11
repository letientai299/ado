package models

// ReferenceLinks represents a collection of REST reference links.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type ReferenceLinks struct {
	// The dictionary of links.
	// The key is the link relationship, and the value is the link object (usually
	// having a href property).
	Links map[string]any `json:"links,omitempty"`
}
