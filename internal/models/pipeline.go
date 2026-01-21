package models

import "time"

// DefinitionQueueStatus indicates whether builds can be queued against a definition.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#definitionqueuestatus
type DefinitionQueueStatus string

const (
	// DefinitionQueueStatusEnabled indicates the definition is enabled and builds can be queued.
	DefinitionQueueStatusEnabled DefinitionQueueStatus = "enabled"
	// DefinitionQueueStatusDisabled indicates the definition is disabled and builds cannot be
	// queued.
	DefinitionQueueStatusDisabled DefinitionQueueStatus = "disabled"
	// DefinitionQueueStatusPaused indicates the definition is paused. Scheduled builds will not
	// queue, but manual builds and CI/PR triggers will still queue.
	DefinitionQueueStatusPaused DefinitionQueueStatus = "paused"
)

// DefinitionType represents the type of a build definition.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#definitiontype
type DefinitionType string

const (
	// DefinitionTypeXaml indicates a XAML build definition (legacy).
	DefinitionTypeXaml DefinitionType = "xaml"
	// DefinitionTypeBuild indicates a standard build definition.
	DefinitionTypeBuild DefinitionType = "build"
)

// BuildStatus represents the status of a build.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#buildstatus
type BuildStatus string

const (
	// BuildStatusNone indicates no status set.
	BuildStatusNone BuildStatus = "none"
	// BuildStatusInProgress indicates the build is currently executing.
	BuildStatusInProgress BuildStatus = "inProgress"
	// BuildStatusCompleted indicates the build has completed.
	BuildStatusCompleted BuildStatus = "completed"
	// BuildStatusCancelling indicates the build is being cancelled.
	BuildStatusCancelling BuildStatus = "cancelling"
	// BuildStatusPostponed indicates the build is postponed.
	BuildStatusPostponed BuildStatus = "postponed"
	// BuildStatusNotStarted indicates the build has not yet started.
	BuildStatusNotStarted BuildStatus = "notStarted"
	// BuildStatusAll represents all statuses (used for filtering).
	BuildStatusAll BuildStatus = "all"
)

// BuildResult represents the final result of a completed build.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#buildresult
type BuildResult string

const (
	// BuildResultNone indicates no result set.
	BuildResultNone BuildResult = "none"
	// BuildResultSucceeded indicates the build completed successfully.
	BuildResultSucceeded BuildResult = "succeeded"
	// BuildResultPartiallySucceeded indicates the build completed with some failures.
	BuildResultPartiallySucceeded BuildResult = "partiallySucceeded"
	// BuildResultFailed indicates the build failed.
	BuildResultFailed BuildResult = "failed"
	// BuildResultCanceled indicates the build was canceled.
	BuildResultCanceled BuildResult = "canceled"
)

// BuildReason represents the reason a build was created.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#buildreason
type BuildReason string

const (
	// BuildReasonNone indicates no reason specified.
	BuildReasonNone BuildReason = "none"
	// BuildReasonManual indicates the build was started manually.
	BuildReasonManual BuildReason = "manual"
	// BuildReasonIndividualCI indicates the build was triggered by a continuous integration
	// trigger for an individual change.
	BuildReasonIndividualCI BuildReason = "individualCI"
	// BuildReasonBatchedCI indicates the build was triggered by a batched continuous integration
	// trigger.
	BuildReasonBatchedCI BuildReason = "batchedCI"
	// BuildReasonSchedule indicates the build was triggered by a schedule.
	BuildReasonSchedule BuildReason = "schedule"
	// BuildReasonScheduleForced indicates the build was triggered by a schedule even though no
	// changes were detected.
	BuildReasonScheduleForced BuildReason = "scheduleForced"
	// BuildReasonUserCreated indicates the build was created by a user.
	BuildReasonUserCreated BuildReason = "userCreated"
	// BuildReasonValidateShelveset indicates the build was created to validate a shelveset.
	BuildReasonValidateShelveset BuildReason = "validateShelveset"
	// BuildReasonCheckInShelveset indicates the build was created to check in a shelveset.
	BuildReasonCheckInShelveset BuildReason = "checkInShelveset"
	// BuildReasonPullRequest indicates the build was triggered by a pull request.
	BuildReasonPullRequest BuildReason = "pullRequest"
	// BuildReasonBuildCompletion indicates the build was triggered by another build completing.
	BuildReasonBuildCompletion BuildReason = "buildCompletion"
	// BuildReasonResourceTrigger indicates the build was triggered by a resource trigger.
	BuildReasonResourceTrigger BuildReason = "resourceTrigger"
	// BuildReasonTriggered indicates the build was triggered (general).
	BuildReasonTriggered BuildReason = "triggered"
	// BuildReasonAll represents all reasons (used for filtering).
	BuildReasonAll BuildReason = "all"
)

// QueuePriority represents the priority of a build in the queue.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#queuepriority
type QueuePriority string

const (
	// QueuePriorityLow indicates low priority.
	QueuePriorityLow QueuePriority = "low"
	// QueuePriorityBelowNormal indicates below normal priority.
	QueuePriorityBelowNormal QueuePriority = "belowNormal"
	// QueuePriorityNormal indicates normal priority (default).
	QueuePriorityNormal QueuePriority = "normal"
	// QueuePriorityAboveNormal indicates above normal priority.
	QueuePriorityAboveNormal QueuePriority = "aboveNormal"
	// QueuePriorityHigh indicates high priority.
	QueuePriorityHigh QueuePriority = "high"
)

// BuildDefinition represents a pipeline definition from the Build API.
// This includes the complete configuration for a CI/CD pipeline including
// triggers, variables, the build process, and queue settings.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#builddefinition
type BuildDefinition struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Id is the unique identifier of the definition.
	Id int32 `json:"id"`

	// Name is the display name of the definition.
	Name string `json:"name"`

	// Path is the folder path where the definition is stored.
	// Uses backslash as separator (e.g., "\\folder\\subfolder").
	Path string `json:"path"`

	// Type indicates the type of the definition (build or xaml).
	Type DefinitionType `json:"type"`

	// QueueStatus indicates whether builds can be queued against this definition.
	QueueStatus DefinitionQueueStatus `json:"queueStatus"`

	// Revision is the definition revision number. Incremented each time the definition is saved.
	Revision int32 `json:"revision"`

	// CreatedDate is when this version of the definition was created.
	CreatedDate *time.Time `json:"createdDate,omitempty"`

	// AuthoredBy is the identity that authored the definition.
	AuthoredBy *IdentityRef `json:"authoredBy,omitempty"`

	// Project is the team project containing this definition.
	Project *TeamProject `json:"project,omitempty"`

	// Queue is the default agent pool queue for builds against this definition.
	Queue *AgentQueue `json:"queue,omitempty"`

	// Process is the build process configuration (YAML or designer).
	Process *Process `json:"process,omitempty"`

	// Repository is the source repository configuration.
	Repository *BuildRepo `json:"repository,omitempty"`

	// Uri is the full URI of the definition.
	Uri string `json:"uri,omitempty"`

	// Url is the REST API URL of the definition.
	Url string `json:"url,omitempty"`

	// Quality indicates the quality of the definition document (draft, etc.).
	Quality string `json:"quality,omitempty"`

	// LatestBuild is the most recent build against this definition.
	LatestBuild *Build `json:"latestBuild,omitempty"`

	// LatestCompletedBuild is the most recent completed build against this definition.
	LatestCompletedBuild *Build `json:"latestCompletedBuild,omitempty"`

	// Description is a description of the definition.
	Description string `json:"description,omitempty"`

	// BadgeEnabled indicates whether badges are enabled for this definition.
	BadgeEnabled bool `json:"badgeEnabled,omitempty"`

	// BuildNumberFormat is the format string used to generate build numbers.
	BuildNumberFormat string `json:"buildNumberFormat,omitempty"`

	// Comment is a save-time comment for the definition revision.
	Comment string `json:"comment,omitempty"`

	// JobAuthorizationScope is the authorization scope for jobs queued against this definition.
	JobAuthorizationScope string `json:"jobAuthorizationScope,omitempty"`

	// JobTimeoutInMinutes is the default execution timeout (in minutes) for builds.
	JobTimeoutInMinutes int32 `json:"jobTimeoutInMinutes,omitempty"`

	// JobCancelTimeoutInMinutes is the timeout (in minutes) for cancelled builds.
	JobCancelTimeoutInMinutes int32 `json:"jobCancelTimeoutInMinutes,omitempty"`

	// Tags are the tags associated with the definition.
	Tags []string `json:"tags,omitempty"`

	// Variables is a dictionary of build definition variables.
	Variables map[string]BuildDefinitionVariable `json:"variables,omitempty"`

	// DraftOf is a reference to the definition this is a draft of, if applicable.
	DraftOf *BuildDefinition `json:"draftOf,omitempty"`

	// Drafts are the drafts associated with this definition.
	Drafts []BuildDefinition `json:"drafts,omitempty"`
}

// BuildDefinitionVariable represents a variable in a build definition.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#builddefinitionvariable
type BuildDefinitionVariable struct {
	// Value is the value of the variable.
	Value string `json:"value,omitempty"`

	// IsSecret indicates whether the variable is a secret.
	// Secret variables are masked in logs and cannot be read back from the API.
	IsSecret bool `json:"isSecret,omitempty"`

	// AllowOverride indicates whether users can override this variable when queuing builds.
	AllowOverride bool `json:"allowOverride,omitempty"`
}

// AgentQueue represents a build queue that builds can be run against.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#agentpoolqueue
type AgentQueue struct {
	// Id is the unique identifier of the queue.
	Id int32 `json:"id"`

	// Name is the display name of the queue.
	Name string `json:"name"`

	// Pool is the agent pool backing this queue.
	Pool *AgentPool `json:"pool,omitempty"`

	// Url is the REST API URL of the queue.
	Url string `json:"url,omitempty"`
}

// AgentPool represents an agent pool containing build agents.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/distributedtask/pools/get#taskagentpool
type AgentPool struct {
	// Id is the unique identifier of the pool.
	Id int32 `json:"id"`

	// Name is the display name of the pool.
	Name string `json:"name"`

	// IsHosted indicates whether this is a Microsoft-hosted pool.
	// Hosted pools are managed by Azure DevOps and provide fresh VMs for each build.
	IsHosted bool `json:"isHosted,omitempty"`

	// PoolType indicates the type of the pool (automation, deployment, etc.).
	PoolType string `json:"poolType,omitempty"`

	// Size is the current size of the pool.
	Size int32 `json:"size,omitempty"`
}

// Process represents the build process configuration.
// This describes how the build is executed (YAML file path or designer process).
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#buildprocess
type Process struct {
	// Type indicates the process type.
	// 1 = Designer process (phase-based)
	// 2 = YAML process
	Type int32 `json:"type"`

	// YamlFilename is the path to the YAML file (only for type=2).
	YamlFilename string `json:"yamlFilename,omitempty"`
}

// BuildRepo represents the repository configuration for a build definition.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/definitions/get#buildrepository
type BuildRepo struct {
	// Id is the unique identifier of the repository.
	// For Azure Repos Git, this is the repository GUID.
	Id string `json:"id"`

	// Name is the display name of the repository.
	Name string `json:"name"`

	// Type is the repository type (TfsGit, GitHub, GitHubEnterprise, Bitbucket, etc.).
	Type string `json:"type"`

	// Url is the URL of the repository.
	Url string `json:"url,omitempty"`

	// DefaultBranch is the default branch for the repository (e.g., "refs/heads/main").
	DefaultBranch string `json:"defaultBranch,omitempty"`

	// RootFolder is the root folder within the repository.
	RootFolder string `json:"rootFolder,omitempty"`

	// Clean indicates whether to clean the working directory before each build.
	Clean string `json:"clean,omitempty"`

	// CheckoutSubmodules indicates whether to checkout submodules.
	CheckoutSubmodules bool `json:"checkoutSubmodules,omitempty"`
}

// Build represents a build run (execution instance of a pipeline).
// This contains the runtime state of a specific build including its status,
// result, timing information, and source details.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#build
type Build struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Id is the unique identifier of the build.
	Id int32 `json:"id"`

	// BuildNumber is the human-readable build number/name (e.g., "20240115.1").
	BuildNumber string `json:"buildNumber,omitempty"`

	// BuildNumberRevision is the revision of the build number within a day.
	BuildNumberRevision int32 `json:"buildNumberRevision,omitempty"`

	// Status is the current status of the build (none, inProgress, completed, etc.).
	Status BuildStatus `json:"status,omitempty"`

	// Result is the final result of the build (succeeded, failed, canceled, etc.).
	// Only set when Status is "completed".
	Result BuildResult `json:"result,omitempty"`

	// QueueTime is when the build was queued.
	QueueTime *time.Time `json:"queueTime,omitempty"`

	// StartTime is when the build started executing.
	StartTime *time.Time `json:"startTime,omitempty"`

	// FinishTime is when the build completed.
	FinishTime *time.Time `json:"finishTime,omitempty"`

	// SourceBranch is the source branch for the build (e.g., "refs/heads/main").
	SourceBranch string `json:"sourceBranch,omitempty"`

	// SourceVersion is the source commit SHA that was built.
	SourceVersion string `json:"sourceVersion,omitempty"`

	// SourceVersionMessage is the commit message of the source version.
	SourceVersionMessage string `json:"sourceVersionMessage,omitempty"`

	// Reason is why the build was created (manual, CI, PR, schedule, etc.).
	Reason BuildReason `json:"reason,omitempty"`

	// RequestedFor is the identity the build was requested on behalf of.
	RequestedFor *IdentityRef `json:"requestedFor,omitempty"`

	// RequestedBy is the identity that queued the build.
	RequestedBy *IdentityRef `json:"requestedBy,omitempty"`

	// TriggeredByBuild is the build that triggered this one (for build completion triggers).
	TriggeredByBuild *Build `json:"triggeredByBuild,omitempty"`

	// TriggerInfo contains source provider-specific trigger information.
	TriggerInfo *TriggerInfo `json:"triggerInfo,omitempty"`

	// Definition is the pipeline definition this build was run against.
	Definition *BuildDefinition `json:"definition,omitempty"`

	// Project is the team project containing this build.
	Project *TeamProject `json:"project,omitempty"`

	// Repository is the repository that was built.
	Repository *BuildRepo `json:"repository,omitempty"`

	// Priority is the build's priority in the queue.
	Priority QueuePriority `json:"priority,omitempty"`

	// Queue is the agent pool queue the build ran against.
	Queue *AgentQueue `json:"queue,omitempty"`

	// QueuePosition is the build's current position in the queue (if queued).
	QueuePosition int32 `json:"queuePosition,omitempty"`

	// Logs is a reference to the build logs.
	Logs *BuildLogReference `json:"logs,omitempty"`

	// Uri is the full URI of the build.
	Uri string `json:"uri,omitempty"`

	// Url is the REST API URL of the build.
	Url string `json:"url,omitempty"`

	// Tags are the tags associated with this build.
	Tags []string `json:"tags,omitempty"`

	// Parameters are the parameters passed to the build (as JSON string).
	Parameters string `json:"parameters,omitempty"`

	// TemplateParameters are the template parameters passed to the YAML pipeline.
	TemplateParameters map[string]string `json:"templateParameters,omitempty"`

	// LastChangedBy is the identity that last modified the build.
	LastChangedBy *IdentityRef `json:"lastChangedBy,omitempty"`

	// LastChangedDate is when the build was last modified.
	LastChangedDate *time.Time `json:"lastChangedDate,omitempty"`

	// Deleted indicates whether the build has been deleted.
	Deleted bool `json:"deleted,omitempty"`

	// DeletedBy is the identity that deleted the build.
	DeletedBy *IdentityRef `json:"deletedBy,omitempty"`

	// DeletedDate is when the build was deleted.
	DeletedDate *time.Time `json:"deletedDate,omitempty"`

	// DeletedReason describes why/how the build was deleted.
	DeletedReason string `json:"deletedReason,omitempty"`

	// RetainedByRelease indicates whether the build is retained by a release.
	RetainedByRelease bool `json:"retainedByRelease,omitempty"`
}

// BuildLogReference represents a reference to the logs for a build.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#buildlogreference
type BuildLogReference struct {
	// Id is the unique identifier of the log.
	Id int32 `json:"id"`

	// Type is the type of the log location.
	Type string `json:"type,omitempty"`

	// Url is the URL to download the logs.
	Url string `json:"url,omitempty"`
}

// TriggerInfo contains information about what triggered a build.
// The fields populated depend on the trigger type (CI, PR, schedule, etc.).
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/builds/get#build
type TriggerInfo struct {
	// CiSourceBranch is the source branch for CI-triggered builds.
	CiSourceBranch string `json:"ci.sourceBranch,omitempty"`

	// CiSourceSha is the commit SHA for CI-triggered builds.
	CiSourceSha string `json:"ci.sourceSha,omitempty"`

	// CiMessage is the commit message for CI-triggered builds.
	CiMessage string `json:"ci.message,omitempty"`

	// PrNumber is the pull request number for PR-triggered builds.
	PrNumber string `json:"pr.number,omitempty"`

	// PrSourceBranch is the source branch of the PR for PR-triggered builds.
	PrSourceBranch string `json:"pr.sourceBranch,omitempty"`

	// PrTargetBranch is the target branch of the PR for PR-triggered builds.
	PrTargetBranch string `json:"pr.targetBranch,omitempty"`

	// PrTitle is the title of the PR for PR-triggered builds.
	PrTitle string `json:"pr.title,omitempty"`

	// PrSourceSha is the source commit SHA of the PR for PR-triggered builds.
	PrSourceSha string `json:"pr.sourceSha,omitempty"`

	// ScheduledReason describes why a scheduled build was triggered.
	ScheduledReason string `json:"scheduledReason,omitempty"`

	// TriggeredBuildId is the ID of the build that triggered this one
	// (for build completion triggers).
	TriggeredBuildId string `json:"triggeredBuildId,omitempty"`
}

// TimelineRecordState represents the state of a timeline record.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#timelinerecordstate
type TimelineRecordState string

const (
	// TimelineRecordStatePending indicates the record is waiting to start.
	TimelineRecordStatePending TimelineRecordState = "pending"
	// TimelineRecordStateInProgress indicates the record is currently executing.
	TimelineRecordStateInProgress TimelineRecordState = "inProgress"
	// TimelineRecordStateCompleted indicates the record has finished.
	TimelineRecordStateCompleted TimelineRecordState = "completed"
)

// TaskResult represents the result of a task or timeline record.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#taskresult
type TaskResult string

const (
	// TaskResultSucceeded indicates the task completed successfully.
	TaskResultSucceeded TaskResult = "succeeded"
	// TaskResultSucceededWithIssues indicates the task completed with warnings.
	TaskResultSucceededWithIssues TaskResult = "succeededWithIssues"
	// TaskResultFailed indicates the task failed.
	TaskResultFailed TaskResult = "failed"
	// TaskResultCanceled indicates the task was canceled.
	TaskResultCanceled TaskResult = "canceled"
	// TaskResultSkipped indicates the task was skipped.
	TaskResultSkipped TaskResult = "skipped"
	// TaskResultAbandoned indicates the task was abandoned.
	TaskResultAbandoned TaskResult = "abandoned"
)

// Timeline represents the build timeline containing all stages, jobs, and tasks.
// The timeline provides a hierarchical view of build execution with parent-child
// relationships between stages, jobs, and tasks.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#timeline
type Timeline struct {
	// Id is the unique identifier of the timeline.
	Id string `json:"id"`

	// ChangeId is used for optimistic concurrency control.
	ChangeId int32 `json:"changeId,omitempty"`

	// LastChangedBy is the identity that last modified the timeline.
	LastChangedBy string `json:"lastChangedBy,omitempty"`

	// LastChangedOn is when the timeline was last modified.
	LastChangedOn *time.Time `json:"lastChangedOn,omitempty"`

	// Records are the timeline entries (stages, jobs, tasks).
	Records []TimelineRecord `json:"records"`

	// Url is the REST API URL of the timeline.
	Url string `json:"url,omitempty"`
}

// TimelineRecord represents a stage, job, or task in the build timeline.
// Records form a hierarchy: Stage -> Job -> Task.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#timelinerecord
type TimelineRecord struct {
	// Links contains REST reference links for related resources.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Id is the unique identifier of the record.
	Id string `json:"id"`

	// ParentId is the ID of the parent record (e.g., job's parent is its stage).
	ParentId string `json:"parentId,omitempty"`

	// Type indicates the record type: "Stage", "Job", "Task", "Phase", or "Checkpoint".
	Type string `json:"type"`

	// Name is the display name of the record.
	Name string `json:"name"`

	// Identifier is a string identifier consistent across retry attempts.
	Identifier string `json:"identifier,omitempty"`

	// State is the current state of the record (pending, inProgress, completed).
	State TimelineRecordState `json:"state"`

	// Result is the final result of the record (succeeded, failed, canceled, skipped).
	// Only set when State is "completed".
	Result TaskResult `json:"result,omitempty"`

	// ResultCode is an optional result code providing additional details.
	ResultCode string `json:"resultCode,omitempty"`

	// Order is the ordinal position relative to sibling records.
	Order int32 `json:"order"`

	// StartTime is when the record started executing.
	StartTime *time.Time `json:"startTime,omitempty"`

	// FinishTime is when the record finished.
	FinishTime *time.Time `json:"finishTime,omitempty"`

	// PercentComplete is the current completion percentage (0-100).
	PercentComplete int32 `json:"percentComplete,omitempty"`

	// CurrentOperation describes what the record is currently doing.
	CurrentOperation string `json:"currentOperation,omitempty"`

	// Log is a reference to the log file for this record.
	Log *LogRef `json:"log,omitempty"`

	// ErrorCount is the number of errors produced by this record.
	ErrorCount int32 `json:"errorCount,omitempty"`

	// WarningCount is the number of warnings produced by this record.
	WarningCount int32 `json:"warningCount,omitempty"`

	// WorkerName is the name of the agent that ran this record.
	WorkerName string `json:"workerName,omitempty"`

	// Attempt is the attempt number (for retries).
	Attempt int32 `json:"attempt,omitempty"`

	// Issues are the errors and warnings associated with this record.
	Issues []Issue `json:"issues,omitempty"`

	// Task is a reference to the task definition (for Task type records).
	Task *TaskReference `json:"task,omitempty"`
}

// LogRef references a log file for a timeline record.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#buildlogreference
type LogRef struct {
	// Id is the unique identifier of the log.
	Id int32 `json:"id"`

	// Type is the type of the log location.
	Type string `json:"type,omitempty"`

	// Url is the URL to download the log content.
	Url string `json:"url"`
}

// Issue represents an error or warning associated with a build timeline record.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#issue
type Issue struct {
	// Type is the issue type: "error" or "warning".
	Type string `json:"type,omitempty"`

	// Category is the issue category.
	Category string `json:"category,omitempty"`

	// Message is the issue message.
	Message string `json:"message,omitempty"`

	// Data contains additional issue data as key-value pairs.
	Data map[string]string `json:"data,omitempty"`
}

// TaskReference is a reference to a task definition.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/build/timeline/get#taskreference
type TaskReference struct {
	// Id is the unique identifier of the task.
	Id string `json:"id,omitempty"`

	// Name is the display name of the task.
	Name string `json:"name,omitempty"`

	// Version is the task version.
	Version string `json:"version,omitempty"`
}
