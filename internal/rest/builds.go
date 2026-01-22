package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
)

// Builds provides access to Azure DevOps Build APIs.
//
// This client wraps the Build REST API for working with build runs,
// timelines, and logs.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds
type Builds struct {
	client Client
}

// ForProject returns a [ProjectBuilds] client scoped to the given repository.
func (b Builds) ForProject(repo config.Repository) ProjectBuilds {
	baseURL, _ := url.JoinPath(
		adoHost,
		repo.Org,
		repo.Project,
		"_apis/build/builds",
	)

	return ProjectBuilds{
		client:  b.client,
		baseURL: baseURL,
		repo:    repo,
	}
}

// ProjectBuilds provides operations on builds within a specific Azure DevOps project.
type ProjectBuilds struct {
	client  Client
	baseURL string
	repo    config.Repository
}

// BuildListOptions configures the builds list query.
type BuildListOptions struct {
	// DefinitionID filters builds for this pipeline definition.
	DefinitionID int32

	// Top limits the maximum number of builds returned.
	Top int

	// StatusFilter filters builds by status (e.g., "completed", "inProgress").
	StatusFilter string

	// ResultFilter filters builds by result (e.g., "succeeded", "failed").
	ResultFilter string

	// BranchName filters builds by source branch (e.g., "refs/heads/main").
	BranchName string
}

// List retrieves builds matching the specified criteria.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/list
func (pb ProjectBuilds) List(ctx context.Context, opts BuildListOptions) ([]models.Build, error) {
	var qs []_shared.Querier

	if opts.DefinitionID > 0 {
		qs = append(qs, _shared.KV[int32]{Key: "definitions", Value: opts.DefinitionID})
	}

	if opts.Top > 0 {
		qs = append(qs, _shared.KV[int]{Key: "$top", Value: opts.Top})
	}

	if opts.StatusFilter != "" {
		qs = append(qs, _shared.KV[string]{Key: "statusFilter", Value: opts.StatusFilter})
	}

	if opts.ResultFilter != "" {
		qs = append(qs, _shared.KV[string]{Key: "resultFilter", Value: opts.ResultFilter})
	}

	if opts.BranchName != "" {
		qs = append(qs, _shared.KV[string]{Key: "branchName", Value: opts.BranchName})
	}

	list, err := httpGet[List[models.Build]](ctx, pb.client, pb.baseURL, qs...)
	if err != nil {
		return nil, err
	}
	return list.Value, nil
}

// ByID retrieves a single build by its ID.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get
func (pb ProjectBuilds) ByID(ctx context.Context, buildID int32) (*models.Build, error) {
	buildURL, _ := url.JoinPath(pb.baseURL, strconv.FormatInt(int64(buildID), 10))
	return httpGet[models.Build](ctx, pb.client, buildURL)
}

// Timeline retrieves the timeline (stages, jobs, tasks) for a build.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get
func (pb ProjectBuilds) Timeline(ctx context.Context, buildID int32) (*models.Timeline, error) {
	timelineURL, _ := url.JoinPath(pb.baseURL, strconv.FormatInt(int64(buildID), 10), "timeline")
	return httpGet[models.Timeline](ctx, pb.client, timelineURL)
}

// LogContent retrieves the text content of a specific log.
// If startLine and endLine are 0, the entire log is returned.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get-build-log
func (pb ProjectBuilds) LogContent(
	ctx context.Context,
	buildID, logID int32,
	startLine, endLine int,
) (string, error) {
	logURL, _ := url.JoinPath(
		pb.baseURL,
		strconv.FormatInt(int64(buildID), 10),
		"logs",
		strconv.FormatInt(int64(logID), 10),
	)

	var qs []_shared.Querier
	if startLine > 0 {
		qs = append(qs, _shared.KV[int]{Key: "startLine", Value: startLine})
	}
	if endLine > 0 {
		qs = append(qs, _shared.KV[int]{Key: "endLine", Value: endLine})
	}

	return httpGetText(ctx, pb.client, logURL, qs...)
}

// httpGetText performs a GET request and returns the response as plain text.
func httpGetText(ctx context.Context, c Client, u string, qs ..._shared.Querier) (string, error) {
	ctx = WithAPIVersion(ctx, apiVersion7_1)
	u = _shared.AppendQueries(u, qs...)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	buf, err := callAndDecode[[]byte](c, req, func(r io.Reader) (*[]byte, error) {
		buf, err := io.ReadAll(r)
		return &buf, err
	})
	if err != nil {
		return "", err
	}

	return string(*buf), nil
}
