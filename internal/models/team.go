package models

// ProjectState represents the state of a team project.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/core/projects/get#projectstate
type ProjectState string

const (
	// ProjectStateDeleting indicates the project is being deleted.
	ProjectStateDeleting ProjectState = "deleting"
	// ProjectStateNew indicates the project is newly created.
	ProjectStateNew ProjectState = "new"
	// ProjectStateWellFormed indicates the project is in a normal state.
	ProjectStateWellFormed ProjectState = "wellFormed"
	// ProjectStateCreatePending indicates the project creation is pending.
	ProjectStateCreatePending ProjectState = "createPending"
	// ProjectStateAll is used to match all project states in queries.
	ProjectStateAll ProjectState = "all"
	// ProjectStateUnchanged indicates the project state is unchanged.
	ProjectStateUnchanged ProjectState = "unchanged"
	// ProjectStateDeleted indicates the project has been deleted.
	ProjectStateDeleted ProjectState = "deleted"
)

// ProjectVisibility represents the visibility level of a team project.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/core/projects/get#projectvisibility
type ProjectVisibility string

const (
	// ProjectVisibilityPrivate indicates the project is only visible to members.
	ProjectVisibilityPrivate ProjectVisibility = "private"
	// ProjectVisibilityPublic indicates the project is publicly visible.
	ProjectVisibilityPublic ProjectVisibility = "public"
)

// TeamProject represents a team project in Azure DevOps.
// A project is a container for source code, work items, builds, and other resources.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/core/projects/get#teamproject
type TeamProject struct {
	// Id is the unique identifier (GUID) of the project.
	Id string `json:"id,omitempty"`

	// Name is the name of the project.
	Name string `json:"name,omitempty"`

	// Description is the description of the project.
	Description string `json:"description,omitempty"`

	// Url is the REST API URL of the project.
	Url string `json:"url,omitempty"`

	// State indicates the current state of the project.
	State ProjectState `json:"state,omitempty"`

	// Revision is the revision number of the project.
	// Incremented when the project is updated.
	Revision int64 `json:"revision,omitempty"`

	// Visibility indicates whether the project is public or private.
	Visibility ProjectVisibility `json:"visibility,omitempty"`

	// LastUpdateTime is when the project was last modified.
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`

	// Abbreviation is the abbreviated name of the project.
	Abbreviation string `json:"abbreviation,omitempty"`

	// DefaultTeamImageUrl is the URL of the default team's image.
	DefaultTeamImageUrl string `json:"defaultTeamImageUrl,omitempty"`
}
