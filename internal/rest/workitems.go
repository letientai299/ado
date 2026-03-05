package rest

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
)

const pathWorkItems = "workitems"

// WorkItems provides access to the Azure DevOps Work Item Tracking REST API.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items
type WorkItems struct {
	client  Client
	org     string
	project string
	baseURL string
}

// WorkItems returns a WorkItems client scoped to the given repository's org and project.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items
func (c Client) WorkItems(repo config.Repository) WorkItems {
	baseURL, _ := url.JoinPath(adoHost, repo.Org, repo.Project, "_apis/wit")
	return WorkItems{
		client:  c,
		org:     repo.Org,
		project: repo.Project,
		baseURL: baseURL,
	}
}

// ByID fetches a single work item by its ID.
// Use expand parameter to include additional data like relations.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item
func (w WorkItems) ByID(
	ctx context.Context,
	id int,
	expand models.WorkItemExpand,
) (*models.WorkItem, error) {
	wiURL, _ := url.JoinPath(w.baseURL, pathWorkItems, strconv.Itoa(id))

	var qs []_shared.Querier
	if expand != "" && expand != models.WorkItemExpandNone {
		qs = append(qs, _shared.KV[string]{Key: "$expand", Value: string(expand)})
	}

	return httpGet[models.WorkItem](ctx, w.client, wiURL, qs...)
}

// List fetches multiple work items by their IDs.
// This is typically used after a WIQL query which only returns IDs.
// Maximum 200 IDs per request.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/list
func (w WorkItems) List(
	ctx context.Context,
	ids []int,
	expand models.WorkItemExpand,
) ([]models.WorkItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	wiURL, _ := url.JoinPath(w.baseURL, pathWorkItems)

	// Convert IDs to comma-separated string
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}

	qs := []_shared.Querier{
		_shared.KV[string]{Key: "ids", Value: strings.Join(idStrs, ",")},
	}

	if expand != "" && expand != models.WorkItemExpandNone {
		qs = append(qs, _shared.KV[string]{Key: "$expand", Value: string(expand)})
	}

	list, err := httpGet[List[models.WorkItem]](ctx, w.client, wiURL, qs...)
	if err != nil {
		return nil, err
	}

	return list.Value, nil
}

// WorkItemDeleteResponse is the response from deleting a work item.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/delete#workitemdelete
type WorkItemDeleteResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Project   string `json:"project"`
	DeletedBy string `json:"deletedBy"`
	DeletedDate string `json:"deletedDate"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	URL       string `json:"url"`
}

// Delete sends a work item to the Recycle Bin (soft delete).
// If destroy is true, the work item is permanently deleted (cannot be recovered).
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/delete
func (w WorkItems) Delete(
	ctx context.Context,
	id int,
	destroy bool,
) (*WorkItemDeleteResponse, error) {
	wiURL, _ := url.JoinPath(w.baseURL, pathWorkItems, strconv.Itoa(id))
	var qs []_shared.Querier
	if destroy {
		qs = append(qs, _shared.KV[string]{Key: "destroy", Value: "true"})
	}
	return httpDelete[WorkItemDeleteResponse](ctx, w.client, wiURL, qs...)
}

// JsonPatchOp represents a single JSON Patch operation.
// Used by the Work Item Create/Update APIs which require Content-Type: application/json-patch+json.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/create#jsonpatchoperation
type JsonPatchOp struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

// Create creates a new work item of the specified type.
// The wiType is the work item type name (e.g., "Bug", "Task", "User Story").
// Fields are set via JSON Patch operations.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/create
func (w WorkItems) Create(
	ctx context.Context,
	wiType string,
	fields []JsonPatchOp,
) (*models.WorkItem, error) {
	wiURL, _ := url.JoinPath(w.baseURL, pathWorkItems, "$"+wiType)
	return httpPatchJsonPatch[models.WorkItem](ctx, w.client, wiURL, fields)
}

// WIQL provides access to the Work Item Query Language API.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql
type WIQL struct {
	client  Client
	org     string
	project string
	baseURL string
}

// WIQL returns a WIQL client scoped to the given repository's org and project.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql
func (c Client) WIQL(repo config.Repository) WIQL {
	baseURL, _ := url.JoinPath(adoHost, repo.Org, repo.Project, "_apis/wit/wiql")
	return WIQL{
		client:  c,
		org:     repo.Org,
		project: repo.Project,
		baseURL: baseURL,
	}
}

// WIQLQuery represents the request body for a WIQL query.
type WIQLQuery struct {
	// The WIQL query string.
	// Example: SELECT [System.Id], [System.Title] FROM WorkItems WHERE [System.State] = 'Active'
	// https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax
	Query string `json:"query"`
}

// Query executes a WIQL query and returns the matching work item references.
// Note: WIQL queries only return work item IDs and URLs. Use WorkItems.List()
// to fetch full work item details.
//
// The top parameter limits the number of results (max 20000, default 200).
// Set top to 0 to use the server default.
//
// Example WIQL:
//
//	SELECT [System.Id], [System.Title], [System.State]
//	FROM WorkItems
//	WHERE [System.WorkItemType] = 'Task'
//	  AND [System.State] <> 'Closed'
//	  AND [System.AssignedTo] = @Me
//	ORDER BY [System.ChangedDate] DESC
//
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql
func (w WIQL) Query(ctx context.Context, query string, top int) (*models.WIQLResult, error) {
	body := WIQLQuery{Query: query}

	var qs []_shared.Querier
	if top > 0 {
		qs = append(qs, _shared.KV[int]{Key: "$top", Value: top})
	}

	wiqlURL := w.baseURL
	if len(qs) > 0 {
		wiqlURL = _shared.AppendQueries(wiqlURL, qs...)
	}

	return httpPost[models.WIQLResult](ctx, w.client, wiqlURL, body)
}
