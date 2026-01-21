package models

import "time"

// PolicyEvaluationRecord represents a policy evaluation for a pull request.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/evaluations/list?view=azure-devops-rest-7.1
type PolicyEvaluationRecord struct {
	// Links to related resources
	Links *ReferenceLinks `json:"_links,omitempty"`

	// ID of the policy evaluation
	EvaluationId string `json:"evaluationId,omitempty"`

	// ID of the policy configuration
	ConfigurationId int `json:"configurationId,omitempty"`

	// The policy configuration
	Configuration PolicyConfiguration `json:"configuration,omitempty"`

	// Context information about the evaluation
	// For build validations, this typically contains build information
	Context map[string]interface{} `json:"context,omitempty"`

	// Status of the policy evaluation: queued, running, approved, rejected, notApplicable, broken
	Status string `json:"status,omitempty"`

	// Whether the policy blocks completion
	Blocking bool `json:"blocking,omitempty"`

	// When the evaluation started
	StartedDate *time.Time `json:"startedDate,omitempty"`

	// When the evaluation finished
	FinishDate *time.Time `json:"finishDate,omitempty"`

	// Artifact ID this evaluation is for
	ArtifactId string `json:"artifactId,omitempty"`
}

// PolicyConfiguration represents a branch policy configuration.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy/configurations?view=azure-devops-rest-7.1
type PolicyConfiguration struct {
	// ID of the configuration
	Id int `json:"id,omitempty"`

	// Type of the policy
	Type PolicyTypeRef `json:"type,omitempty"`

	// Settings for the policy
	Settings map[string]interface{} `json:"settings,omitempty"`

	// Whether the policy is enabled
	IsEnabled bool `json:"isEnabled,omitempty"`

	// Whether the policy is blocking
	IsBlocking bool `json:"isBlocking,omitempty"`
}

// PolicyTypeRef represents a reference to a policy type.
type PolicyTypeRef struct {
	// ID of the policy type.
	// Build validation policy type ID: "0609b952-1397-4640-95ec-e00a01b2c241"
	Id string `json:"id,omitempty"`

	// Display name of the policy type
	DisplayName string `json:"displayName,omitempty"`

	// URL to the policy type
	Url string `json:"url,omitempty"`
}

// Build validation policy type ID constant
const PolicyTypeBuildValidation = "0609b952-1397-4640-95ec-e00a01b2c241"
