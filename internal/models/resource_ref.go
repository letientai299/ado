package models

// ResourceRef represents a reference to a resource.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type ResourceRef struct {
	Id  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}
