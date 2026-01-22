package models

import "time"

// PolicyEvaluationStatus represents the status of a policy evaluation.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/list#policyevaluationstatus
type PolicyEvaluationStatus string

const (
	// PolicyEvaluationStatusQueued indicates the policy is queued to run,
	// or is waiting for some event before progressing.
	PolicyEvaluationStatusQueued PolicyEvaluationStatus = "queued"

	// PolicyEvaluationStatusRunning indicates the policy is currently running.
	PolicyEvaluationStatusRunning PolicyEvaluationStatus = "running"

	// PolicyEvaluationStatusApproved indicates the policy has been fulfilled
	// for this pull request.
	PolicyEvaluationStatusApproved PolicyEvaluationStatus = "approved"

	// PolicyEvaluationStatusRejected indicates the policy has rejected this
	// pull request.
	PolicyEvaluationStatusRejected PolicyEvaluationStatus = "rejected"

	// PolicyEvaluationStatusNotApplicable indicates the policy does not apply
	// to this pull request.
	PolicyEvaluationStatusNotApplicable PolicyEvaluationStatus = "notApplicable"

	// PolicyEvaluationStatusBroken indicates the policy has encountered an
	// unexpected error.
	PolicyEvaluationStatusBroken PolicyEvaluationStatus = "broken"
)

// PolicyEvaluationRecord represents a policy evaluation for a pull request.
// Each record tracks the evaluation of a single policy against a specific
// artifact (typically a pull request).
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/list#policyevaluationrecord
type PolicyEvaluationRecord struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// EvaluationId is a GUID which uniquely identifies this evaluation record.
	// Each policy running on each pull request gets a unique evaluation ID.
	EvaluationId string `json:"evaluationId,omitempty"`

	// ConfigurationId is the ID of the policy configuration.
	ConfigurationId int `json:"configurationId,omitempty"`

	// Configuration contains all configuration data for the policy being evaluated.
	Configuration PolicyConfiguration `json:"configuration,omitempty"`

	// Context contains internal context data of this policy evaluation.
	// For build validation policies, this typically contains build information
	// including the build ID and definition ID.
	Context map[string]any `json:"context,omitempty"`

	// Status is the current status of the policy evaluation
	// (queued, running, approved, rejected, notApplicable, broken).
	Status PolicyEvaluationStatus `json:"status,omitempty"`

	// Blocking indicates whether the policy blocks pull request completion
	// when not approved.
	Blocking bool `json:"blocking,omitempty"`

	// StartedDate is when the policy was first evaluated on this pull request.
	StartedDate *time.Time `json:"startedDate,omitempty"`

	// CompletedDate is when the policy finished evaluating on this pull request.
	// This field was previously named FinishDate in some API versions.
	CompletedDate *time.Time `json:"completedDate,omitempty"`

	// ArtifactId uniquely identifies the target of the policy evaluation.
	// Format: "vstfs:///CodeReview/CodeReviewId/{projectId}/{pullRequestId}"
	ArtifactId string `json:"artifactId,omitempty"`
}

// PolicyConfiguration represents a branch policy configuration.
// Policies enforce rules on pull requests, such as requiring builds to pass,
// minimum reviewer counts, or work item linking.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/get#policyconfiguration
type PolicyConfiguration struct {
	// Id is the unique identifier of the configuration.
	Id int `json:"id,omitempty"`

	// Type is a reference to the policy type definition.
	Type PolicyTypeRef `json:"type,omitempty"`

	// Revision is the policy configuration revision number.
	// Incremented each time the configuration is updated.
	Revision int `json:"revision,omitempty"`

	// Settings contains policy-specific settings as key-value pairs.
	// The structure depends on the policy type. For build validation policies,
	// this includes buildDefinitionId, displayName, validDuration, etc.
	Settings map[string]any `json:"settings,omitempty"`

	// IsEnabled indicates whether the policy is currently enabled.
	IsEnabled bool `json:"isEnabled,omitempty"`

	// IsBlocking indicates whether the policy blocks pull request completion
	// when not satisfied. Non-blocking policies show as warnings.
	IsBlocking bool `json:"isBlocking,omitempty"`

	// IsDeleted indicates whether the policy configuration has been deleted.
	IsDeleted bool `json:"isDeleted,omitempty"`

	// IsEnterpriseManaged indicates whether the policy is managed at the
	// enterprise/organization level rather than the project level.
	IsEnterpriseManaged bool `json:"isEnterpriseManaged,omitempty"`

	// Url is the REST API URL of the policy configuration.
	Url string `json:"url,omitempty"`
}

// PolicyTypeRef represents a reference to a policy type definition.
// Policy types define the behavior and settings schema for policies.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/types/get#policytype
type PolicyTypeRef struct {
	// Id is the unique identifier (GUID) of the policy type.
	// Well-known policy type IDs:
	//   - Build validation: "0609b952-1397-4640-95ec-e00a01b2c241"
	//   - Minimum reviewers: "fa4e907d-c16b-4a4c-9dfa-4906e5d171dd"
	//   - Required reviewers: "fd2167ab-b0be-447a-8ec8-39368250530e"
	//   - Work item linking: "40e92b44-2fe1-4dd6-b3d8-74a9c21d0c6e"
	//   - Comment requirements: "c6a1889d-b943-4856-b76f-9e46bb6b0df2"
	Id string `json:"id,omitempty"`

	// DisplayName is the human-readable name of the policy type.
	DisplayName string `json:"displayName,omitempty"`

	// Url is the REST API URL of the policy type.
	Url string `json:"url,omitempty"`
}

// Well-known policy type IDs for Azure DevOps branch policies.
// See: https://learn.microsoft.com/en-us/azure/devops/repos/git/branch-policies
//
//goland:noinspection GoUnusedConst
const (
	// PolicyTypeBuildValidation is the policy type ID for build validation policies.
	// These policies require a successful build before a pull request can be completed.
	PolicyTypeBuildValidation = "0609b952-1397-4640-95ec-e00a01b2c241"

	// PolicyTypeMinimumReviewers is the policy type ID for minimum reviewer policies.
	// These policies require a minimum number of approvals before completion.
	PolicyTypeMinimumReviewers = "fa4e907d-c16b-4a4c-9dfa-4906e5d171dd"

	// PolicyTypeRequiredReviewers is the policy type ID for required reviewer policies.
	// These policies require specific users or groups to approve.
	PolicyTypeRequiredReviewers = "fd2167ab-b0be-447a-8ec8-39368250530e"

	// PolicyTypeWorkItemLinking is the policy type ID for work item linking policies.
	// These policies require linked work items before completion.
	PolicyTypeWorkItemLinking = "40e92b44-2fe1-4dd6-b3d8-74a9c21d0c6e"

	// PolicyTypeCommentRequirements is the policy type ID for comment resolution policies.
	// These policies require all comments to be resolved before completion.
	PolicyTypeCommentRequirements = "c6a1889d-b943-4856-b76f-9e46bb6b0df2"
)
