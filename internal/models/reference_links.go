package models

// ReferenceLinks represents a collection of REST reference links.
// These links provide navigation to related resources in the Azure DevOps API.
//
// The links dictionary contains key-value pairs where the key is the link
// relationship name (e.g., "self", "web", "repository") and the value is
// typically an object with an "href" property containing the URL.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get#referencelinks
type ReferenceLinks struct {
	// Links is the dictionary of reference links.
	// Common keys include:
	//   - "self": REST API URL of this resource
	//   - "web": Web browser URL for this resource
	//   - "repository": Link to the containing repository
	//   - "commits": Link to commit history
	//   - "workItems": Link to related work items
	Links map[string]any `json:"links,omitempty"`
}
