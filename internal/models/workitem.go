package models

import "time"

// WorkItem represents an Azure DevOps work item.
// A work item is a unit of work that can be tracked through the Azure DevOps
// project lifecycle (e.g., bugs, tasks, user stories, features, epics).
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item
type WorkItem struct {
	// The work item ID. This is a unique identifier within the project.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitem
	ID int `json:"id"`
	// The revision number of the work item. Increments with each update.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitem
	Rev int `json:"rev"`
	// Map of field reference names to values. Standard fields include:
	//   - System.Title: Work item title (string)
	//   - System.Description: Detailed description (HTML string)
	//   - System.State: Current state like "New", "Active", "Closed" (string)
	//   - System.WorkItemType: Type like "Bug", "Task", "User Story" (string)
	//   - System.AssignedTo: Person assigned to the work item (IdentityRef)
	//   - System.AreaPath: Area path for categorization (string)
	//   - System.IterationPath: Sprint/iteration path (string)
	//   - System.Tags: Semicolon-delimited tags (string)
	//   - System.CreatedDate: When created (datetime)
	//   - System.ChangedDate: When last modified (datetime)
	//   - System.CreatedBy: Who created it (IdentityRef)
	//   - System.ChangedBy: Who last modified it (IdentityRef)
	//   - Microsoft.VSTS.Common.Priority: Priority 1-4 (integer)
	//   - System.Parent: Parent work item ID (integer)
	// https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field
	Fields map[string]any `json:"fields"`
	// Relations to other work items, commits, builds, etc.
	// Only populated when $expand=relations is specified.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitemrelation
	Relations []WorkItemRelation `json:"relations,omitempty"`
	// Full URL to the work item REST API endpoint.
	URL string `json:"url"`
	// Links to related resources like HTML (web UI), workItemType, fields, etc.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#referencelinks
	Links *ReferenceLinks `json:"_links,omitempty"`
	// Link to comment pages for this work item. Only present when comments exist.
	CommentVersionRef *WorkItemCommentVersionRef `json:"commentVersionRef,omitempty"`
}

// WorkItemRelation represents a link between a work item and another artifact.
// Relations can link to other work items (parent, child, related), commits,
// pull requests, builds, or external URLs.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitemrelation
type WorkItemRelation struct {
	// The relation type reference name. Common types include:
	//   - System.LinkTypes.Hierarchy-Forward: Parent link
	//   - System.LinkTypes.Hierarchy-Reverse: Child link
	//   - System.LinkTypes.Related: Related link
	//   - ArtifactLink: Link to commits, PRs, builds, etc.
	// https://learn.microsoft.com/en-us/azure/devops/boards/queries/link-type-reference
	Rel string `json:"rel"`
	// URL to the related artifact or work item.
	URL string `json:"url"`
	// Additional attributes about the relation.
	// For artifact links, may include "name" (display name) and "comment".
	Attributes map[string]any `json:"attributes,omitempty"`
}

// WorkItemCommentVersionRef contains information about comments on a work item.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitemcommentversionref
type WorkItemCommentVersionRef struct {
	// The ID of the comment.
	CommentID int `json:"commentId"`
	// The version of the comment.
	Version int `json:"version"`
	// URL to fetch the comment.
	URL string `json:"url"`
}

// WIQLResult represents the response from a WIQL query.
// WIQL (Work Item Query Language) is a SQL-like language for querying work items.
// Note: WIQL queries only return work item IDs and URLs, not full details.
// Use the returned IDs with the List Work Items API to fetch full details.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql
type WIQLResult struct {
	// The type of query: "flat" for simple queries, "oneHop" or "tree" for
	// link queries that return hierarchical results.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#querytype
	QueryType string `json:"queryType"`
	// The point in time when the query was executed. Useful for understanding
	// when the snapshot of work items was taken.
	AsOf time.Time `json:"asOf"`
	// Columns returned by the query (based on SELECT clause).
	// Each column has ReferenceName (e.g., "System.Id") and Name (e.g., "ID").
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemfieldreference
	Columns []WorkItemFieldReference `json:"columns,omitempty"`
	// List of work item references for flat queries.
	// Only contains ID and URL - use List API to get full details.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemreference
	WorkItems []WorkItemReference `json:"workItems,omitempty"`
	// List of work item link relations for tree/oneHop queries.
	// Contains source and target work item references.
	// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemlink
	WorkItemRelations []WorkItemLink `json:"workItemRelations,omitempty"`
	// Sorting columns for the query results.
	SortColumns []WorkItemQuerySortColumn `json:"sortColumns,omitempty"`
}

// WorkItemReference is a lightweight reference to a work item.
// Contains only the ID and URL, not the full work item details.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemreference
type WorkItemReference struct {
	// The work item ID.
	ID int `json:"id"`
	// URL to fetch the full work item details.
	URL string `json:"url"`
}

// WorkItemLink represents a link relationship in WIQL tree/oneHop queries.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemlink
type WorkItemLink struct {
	// The relation type (e.g., "System.LinkTypes.Hierarchy-Forward").
	Rel string `json:"rel,omitempty"`
	// The source work item reference (may be nil for root items).
	Source *WorkItemReference `json:"source,omitempty"`
	// The target work item reference.
	Target *WorkItemReference `json:"target,omitempty"`
}

// WorkItemFieldReference describes a field in a work item query result.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemfieldreference
type WorkItemFieldReference struct {
	// The reference name of the field (e.g., "System.Id", "System.Title").
	// This is the unique identifier used in WIQL queries and API calls.
	ReferenceName string `json:"referenceName"`
	// The display name of the field (e.g., "ID", "Title").
	Name string `json:"name"`
	// The URL to the field definition.
	URL string `json:"url,omitempty"`
}

// WorkItemQuerySortColumn describes how query results are sorted.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql#workitemquerysortcolumn
type WorkItemQuerySortColumn struct {
	// The field to sort by.
	Field *WorkItemFieldReference `json:"field"`
	// True if sorting in descending order, false for ascending.
	Descending bool `json:"descending"`
}

// WorkItemExpand specifies which additional data to include when fetching work items.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item#workitemexpand
type WorkItemExpand string

const (
	// WorkItemExpandNone returns only basic work item fields.
	WorkItemExpandNone WorkItemExpand = "none"
	// WorkItemExpandRelations includes all relations (links to other work items, commits, etc.).
	WorkItemExpandRelations WorkItemExpand = "relations"
	// WorkItemExpandFields includes field definitions along with values.
	WorkItemExpandFields WorkItemExpand = "fields"
	// WorkItemExpandLinks includes the _links property with REST API links.
	WorkItemExpandLinks WorkItemExpand = "links"
	// WorkItemExpandAll includes relations, fields, and links.
	WorkItemExpandAll WorkItemExpand = "all"
)

// Common work item field reference names.
// These are the standard Azure DevOps field names used in WIQL queries and API responses.
// https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field
const (
	// FieldID is the unique identifier for a work item (integer).
	FieldID = "System.Id"
	// FieldTitle is the work item title (string, max 255 chars).
	FieldTitle = "System.Title"
	// FieldDescription is the detailed description (HTML string).
	FieldDescription = "System.Description"
	// FieldState is the current state (string, e.g., "New", "Active", "Closed").
	FieldState = "System.State"
	// FieldReason is the reason for the current state (string).
	FieldReason = "System.Reason"
	// FieldWorkItemType is the type (string, e.g., "Bug", "Task", "User Story").
	FieldWorkItemType = "System.WorkItemType"
	// FieldAssignedTo is the person assigned to the work item (IdentityRef).
	FieldAssignedTo = "System.AssignedTo"
	// FieldAreaPath is the area path for categorization (string).
	FieldAreaPath = "System.AreaPath"
	// FieldIterationPath is the sprint/iteration path (string).
	FieldIterationPath = "System.IterationPath"
	// FieldTags is semicolon-delimited tags (string, e.g., "tag1; tag2").
	FieldTags = "System.Tags"
	// FieldCreatedDate is when the work item was created (datetime).
	FieldCreatedDate = "System.CreatedDate"
	// FieldCreatedBy is who created the work item (IdentityRef).
	FieldCreatedBy = "System.CreatedBy"
	// FieldChangedDate is when the work item was last modified (datetime).
	FieldChangedDate = "System.ChangedDate"
	// FieldChangedBy is who last modified the work item (IdentityRef).
	FieldChangedBy = "System.ChangedBy"
	// FieldPriority is the priority level (integer, typically 1-4).
	FieldPriority = "Microsoft.VSTS.Common.Priority"
	// FieldParent is the ID of the parent work item (integer).
	FieldParent = "System.Parent"
	// FieldCommentCount is the number of comments on the work item (integer).
	FieldCommentCount = "System.CommentCount"
	// FieldTeamProject is the team project name (string).
	FieldTeamProject = "System.TeamProject"
	// FieldHistory is used to add discussion comments (HTML string, write-only).
	FieldHistory = "System.History"
	// FieldAcceptanceCriteria is the acceptance criteria (HTML string).
	FieldAcceptanceCriteria = "Microsoft.VSTS.Common.AcceptanceCriteria"
	// FieldReproSteps is the reproduction steps for bugs (HTML string).
	FieldReproSteps = "Microsoft.VSTS.TCM.ReproSteps"
)
