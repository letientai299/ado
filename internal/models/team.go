package models

// TeamProject represents a team project.
type TeamProject struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	State       string `json:"state,omitempty"`
	Revision    int64  `json:"revision,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}
