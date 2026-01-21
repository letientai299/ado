package models

import "time"

type DefinitionQueueStatus string

const (
	DefinitionQueueStatusEnabled  DefinitionQueueStatus = "enabled"
	DefinitionQueueStatusDisabled DefinitionQueueStatus = "disabled"
	DefinitionQueueStatusPaused   DefinitionQueueStatus = "paused"
)

// BuildDefinition represents a pipeline definition from the Build API.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/list
type BuildDefinition struct {
	Id                   int32                 `json:"id"`
	Name                 string                `json:"name"`
	Path                 string                `json:"path"`
	Type                 string                `json:"type"`
	QueueStatus          DefinitionQueueStatus `json:"queueStatus"`
	Revision             int32                 `json:"revision"`
	CreatedDate          *time.Time            `json:"createdDate,omitempty"`
	AuthoredBy           *IdentityRef          `json:"authoredBy,omitempty"`
	Project              *TeamProject          `json:"project,omitempty"`
	Queue                *AgentQueue           `json:"queue,omitempty"`
	Process              *Process              `json:"process,omitempty"`
	Repository           *BuildRepo            `json:"repository,omitempty"`
	Uri                  string                `json:"uri,omitempty"`
	Url                  string                `json:"url,omitempty"`
	Quality              string                `json:"quality,omitempty"`
	LatestBuild          *Build                `json:"latestBuild,omitempty"`
	LatestCompletedBuild *Build                `json:"latestCompletedBuild,omitempty"`
}

// AgentQueue represents a build queue.
type AgentQueue struct {
	Id   int32      `json:"id"`
	Name string     `json:"name"`
	Pool *AgentPool `json:"pool,omitempty"`
}

// AgentPool represents an agent pool.
type AgentPool struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	IsHosted bool   `json:"isHosted,omitempty"`
}

// Process represents the build process configuration.
type Process struct {
	Type         int32  `json:"type"`
	YamlFilename string `json:"yamlFilename,omitempty"`
}

// BuildRepo represents the repository configuration for a build.
type BuildRepo struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Url           string `json:"url,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
	RootFolder    string `json:"rootFolder,omitempty"`
}

// Build represents a build run.
type Build struct {
	Id                   int32            `json:"id"`
	BuildNumber          string           `json:"buildNumber,omitempty"`
	BuildNumberRevision  int32            `json:"buildNumberRevision,omitempty"`
	Status               string           `json:"status,omitempty"`
	Result               string           `json:"result,omitempty"`
	QueueTime            *time.Time       `json:"queueTime,omitempty"`
	StartTime            *time.Time       `json:"startTime,omitempty"`
	FinishTime           *time.Time       `json:"finishTime,omitempty"`
	SourceBranch         string           `json:"sourceBranch,omitempty"`
	SourceVersion        string           `json:"sourceVersion,omitempty"`
	SourceVersionMessage string           `json:"sourceVersionMessage,omitempty"`
	Reason               string           `json:"reason,omitempty"`
	RequestedFor         *IdentityRef     `json:"requestedFor,omitempty"`
	RequestedBy          *IdentityRef     `json:"requestedBy,omitempty"`
	TriggeredByBuild     *Build           `json:"triggeredByBuild,omitempty"`
	TriggerInfo          *TriggerInfo     `json:"triggerInfo,omitempty"`
	Definition           *BuildDefinition `json:"definition,omitempty"`
}

// TriggerInfo contains info about what triggered the build.
type TriggerInfo struct {
	CiSourceBranch   string `json:"ci.sourceBranch,omitempty"`
	CiSourceSha      string `json:"ci.sourceSha,omitempty"`
	CiMessage        string `json:"ci.message,omitempty"`
	PrNumber         string `json:"pr.number,omitempty"`
	PrSourceBranch   string `json:"pr.sourceBranch,omitempty"`
	PrTargetBranch   string `json:"pr.targetBranch,omitempty"`
	PrTitle          string `json:"pr.title,omitempty"`
	ScheduledReason  string `json:"scheduledReason,omitempty"`
	TriggeredBuildId string `json:"triggeredBuildId,omitempty"`
}

// Timeline represents the build timeline containing all stages, jobs, and tasks.
type Timeline struct {
	Id      string           `json:"id"`
	Records []TimelineRecord `json:"records"`
}

// TimelineRecord represents a stage, job, or task in the build timeline.
type TimelineRecord struct {
	Id       string  `json:"id"`
	ParentId string  `json:"parentId,omitempty"`
	Type     string  `json:"type"` // "Stage", "Job", "Task"
	Name     string  `json:"name"`
	State    string  `json:"state"` // "pending", "inProgress", "completed"
	Result   string  `json:"result,omitempty"`
	Order    int32   `json:"order"`
	Log      *LogRef `json:"log,omitempty"`
}

// LogRef references a log file for a timeline record.
type LogRef struct {
	Id  int32  `json:"id"`
	Url string `json:"url"`
}
