package rest

import (
	"context"
	"io"
	"net/url"
	"strconv"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
)

// Pipelines provide access to Azure DevOps Build/Pipeline APIs.
//
// This client wraps the Build Definitions REST API, which manages pipeline
// definitions (also known as build definitions) in Azure DevOps. Pipeline
// definitions describe the build process, including triggers, variables,
// and the steps to execute.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions
type Pipelines struct {
	client Client
}

// Definitions return a [PipelineDefinitions] client scoped to the given repository.
// The returned client can be used to list and retrieve pipeline definitions
// associated with the specified repository.
func (p Pipelines) Definitions(repo config.Repository) PipelineDefinitions {
	baseURL, _ := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/build/definitions",
	)

	return PipelineDefinitions{
		client:  p.client,
		baseURL: baseURL,
		repo:    repo,
	}
}

// PipelineDefinitions provides operations on pipeline definitions within
// a specific Azure DevOps project. It is scoped to a repository and provides
// methods to list and retrieve build/pipeline definitions.
//
// Pipeline definitions (also called build definitions) contain the configuration
// for CI/CD pipelines, including the YAML file path, triggers, variables,
// and queue settings.
type PipelineDefinitions struct {
	client  Client
	baseURL string
	repo    config.Repository
}

// List retrieves pipeline definitions matching the specified criteria.
//
// When [ListOptions.RepositoryID] is set, results are filtered to definitions
// that use the specified repository (with TfsGit as the repository type).
//
// The API supports additional query parameters not exposed in [ListOptions]:
//   - builtAfter/notBuiltAfter: filter by build date
//   - definitionIds: retrieve specific definitions by ID
//   - queryOrder: control result ordering
//   - continuationToken: for paginating large result sets
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/list
func (pd PipelineDefinitions) List(
	ctx context.Context,
	opts ListOptions,
) ([]models.BuildDefinition, error) {
	list, err := httpGet[List[models.BuildDefinition]](ctx, pd.client, pd.baseURL, opts)
	if err != nil {
		return nil, err
	}
	return list.Value, nil
}

// ByID retrieves a single pipeline definition by its numeric ID.
//
// The returned [models.BuildDefinition] contains the complete definition including
//   - Repository information and YAML file path
//   - Build process configuration
//   - Triggers, variables, and retention policies
//   - Queue status (enabled/disabled/paused)
//
// The API supports additional query parameters not currently exposed:
//   - revision: retrieve a specific revision instead of latest
//   - includeLatestBuilds: include latest build information
//   - propertyFilters: limit which properties are returned
//
// Returns an error if the definition does not exist or is not accessible.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get
func (pd PipelineDefinitions) ByID(ctx context.Context, id int32) (*models.BuildDefinition, error) {
	defURL, _ := url.JoinPath(pd.baseURL, strconv.FormatInt(int64(id), 10))
	return httpGet[models.BuildDefinition](ctx, pd.client, defURL)
}

// ListOptions configures the pipeline definitions list query.
//
// All fields are optional. When a field is empty/zero, the corresponding
// query parameter is omitted from the API request.
type ListOptions struct {
	// Name filters definitions whose names match this pattern.
	// Supports wildcards (e.g., "build-*" matches "build-main", "build-dev").
	Name string

	// Path filters definitions under this folder path.
	// Use backslash as a separator (e.g., "\\folder\\subfolder").
	// Empty string or "\\" returns definitions at the root.
	Path string

	// RepositoryID filters definitions that use this repository.
	// When set, the query automatically includes repositoryType=TfsGit.
	// This is the GUID of the repository, not the repository name.
	RepositoryID string

	// Top limits the maximum number of definitions returned.
	// When zero, the API default limit applies.
	Top int

	IncludeAllProperties bool
	IncludeLatestBuilds  bool
}

func (l ListOptions) AppendTo(w io.Writer) {
	if l.Name != "" {
		_, _ = w.Write([]byte("&name=" + url.QueryEscape(l.Name)))
	}
	if l.Path != "" {
		_, _ = w.Write([]byte("&path=" + url.QueryEscape(l.Path)))
	}
	if l.RepositoryID != "" {
		_, _ = w.Write([]byte("&repositoryId=" + url.QueryEscape(l.RepositoryID)))

		// TODO (tai): check the repo for its type
		_, _ = w.Write([]byte("&repositoryType=TfsGit"))
	}

	if l.Top > 0 {
		_, _ = w.Write([]byte("&$top=" + strconv.Itoa(l.Top)))
	}

	if l.IncludeAllProperties {
		_, _ = w.Write([]byte("&includeAllProperties=true"))
	}

	if l.IncludeLatestBuilds {
		_, _ = w.Write([]byte("&includeLatestBuilds=true"))
	}
}
