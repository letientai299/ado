package pull_request

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
)

const ErrEmptyTitle util.StrErr = "PR title cannot be empty. PR creation/update cancelled."

type Identity struct {
	Id    string
	Name  string
	Email string
}

type PR struct {
	PullRequestId    int32        `yaml:"pull_request_id" json:"pull_request_id,omitempty"`
	Title            string       `yaml:"title"           json:"title,omitempty"`
	Description      string       `yaml:"description"     json:"description,omitempty"`
	IsDraft          bool         `yaml:"is_draft"        json:"is_draft,omitempty"`
	CreatedBy        Identity     `yaml:"created_by"      json:"created_by"`
	CreationDate     string       `yaml:"creation_date"   json:"creation_date,omitempty"`
	WebURL           string       `yaml:"web_url"         json:"web_url,omitempty"`
	Approvers        []Identity   `yaml:"approvers"       json:"approvers,omitempty"`
	SourceBranchName string       `yaml:"source_branch"   json:"source_branch,omitempty"`
	TargetBranchName string       `yaml:"target_branch"   json:"target_branch,omitempty"`
	BuildStatus      *BuildStatus `yaml:"build_status"    json:"build_status,omitempty"`
	PolicyChecks     PolicyChecks `yaml:"policy_checks"   json:"policy_checks,omitempty"`
}

// PRBuildStatus represents the possible states of a PR build
type PRBuildStatus string

const (
	PRBuildStatusSucceeded PRBuildStatus = "succeeded"
	PRBuildStatusFailed    PRBuildStatus = "failed"
	PRBuildStatusRunning   PRBuildStatus = "running"
	PRBuildStatusPending   PRBuildStatus = "pending"
	PRBuildStatusExpired   PRBuildStatus = "expired"
	PRBuildStatusError     PRBuildStatus = "error"
	PRBuildStatusUnknown   PRBuildStatus = "unknown"
)

// BuildStatus represents the build/pipeline status for a pull request
type BuildStatus struct {
	State       PRBuildStatus `yaml:"state"       json:"state,omitempty"`       // succeeded, failed, pending, expired, error, etc.
	Description string        `yaml:"description" json:"description,omitempty"` // human-readable description
	TargetURL   string        `yaml:"target_url"  json:"target_url,omitempty"`  // link to build results
	Icon        string        `yaml:"icon"        json:"icon,omitempty"`        // emoji icon
	StatusText  string        `yaml:"status_text" json:"status_text,omitempty"` // display text for the status
}

// PolicyCheck represents a simplified policy check result
type PolicyCheck struct {
	Name       string `yaml:"name"        json:"name,omitempty"`        // Display the name of the policy
	Status     string `yaml:"status"      json:"status,omitempty"`      // approved, rejected, running, queued, etc.
	IsRequired bool   `yaml:"is_required" json:"is_required,omitempty"` // Whether this is a required check
	Icon       string `yaml:"icon"        json:"icon,omitempty"`        // Status icon for display
}

// PolicyChecks is a slice of PolicyCheck with helper methods
type PolicyChecks []PolicyCheck

// Pending returns only the pending or failed required policies
func (pc PolicyChecks) Pending() PolicyChecks {
	var pending PolicyChecks
	for _, check := range pc {
		if check.IsRequired && (check.isPending() || check.isFailed()) {
			pending = append(pending, check)
		}
	}
	return pending
}

// SummaryText returns a formatted summary of failed/pending policies
func (pc PolicyChecks) SummaryText() string {
	// Pre-allocate with estimated size to reduce allocations
	failedChecks := make([]string, 0, len(pc))
	pendingChecks := make([]string, 0, len(pc))

	for _, check := range pc {
		if !check.IsRequired {
			continue
		}
		if check.isFailed() {
			failedChecks = append(failedChecks, check.Name)
		} else if check.isPending() {
			pendingChecks = append(pendingChecks, check.Name)
		}
	}

	if len(failedChecks) > 0 {
		if len(failedChecks) == 1 {
			return "Failed: " + failedChecks[0]
		}
		return fmt.Sprintf(
			"%d checks failed: %s",
			len(failedChecks),
			strings.Join(failedChecks, ", "),
		)
	}

	if len(pendingChecks) > 0 {
		if len(pendingChecks) == 1 {
			return "Pending: " + pendingChecks[0]
		}
		return fmt.Sprintf(
			"%d checks pending: %s",
			len(pendingChecks),
			strings.Join(pendingChecks, ", "),
		)
	}

	return ""
}

// SummaryIcon returns the appropriate icon for the policy status
func (pc PolicyChecks) SummaryIcon() string {
	for _, check := range pc {
		if check.IsRequired && check.isFailed() {
			return styles.Error(ui.IconFailure)
		}
	}
	for _, check := range pc {
		if check.IsRequired && check.isPending() {
			return styles.Pending(ui.IconPending)
		}
	}
	return ""
}

// isPending checks if a policy is pending or running
func (c PolicyCheck) isPending() bool {
	return c.Status == "queued" || c.Status == "running"
}

// isFailed checks if a policy has failed
func (c PolicyCheck) isFailed() bool {
	return c.Status == "rejected" || c.Status == "broken"
}

func isDirectlyApproved(vote *models.IdentityRefWithVote) bool {
	// NOTE (tai): there's Approve (vote=10) and Approve-with-suggestions (vote=5).
	//  Almost no one use approve with suggestions, so we only consider vote > 0 as approval.
	//  We might update this logic if there's a need later
	return vote.Vote > 0 && !vote.IsContainer
}

func toIdentity(a *models.IdentityRef) Identity {
	return Identity{
		Id:    a.Id,
		Name:  a.DisplayName,
		Email: a.UniqueName,
	}
}

func webURL(baseURL string, id int32) string {
	return baseURL + "/" + strconv.FormatInt(int64(id), 10)
}

func converterWithStatuses(
	baseURL string,
	orgName string,
	repo *models.GitRepository,
	evaluations map[int32][]models.PolicyEvaluationRecord,
) func(m models.GitPullRequest) PR {
	return func(m models.GitPullRequest) PR {
		pr := PR{
			PullRequestId:    m.PullRequestId,
			Title:            m.Title,
			Description:      m.Description,
			IsDraft:          m.IsDraft,
			SourceBranchName: cleanBranchName(m.SourceRefName),
			TargetBranchName: cleanBranchName(m.TargetRefName),
		}

		if m.CreationDate != nil {
			pr.CreationDate = m.CreationDate.Format("2006-01-02")
		}

		pr.WebURL = webURL(baseURL, pr.PullRequestId)

		// Optimize approver extraction - single pass filtering + mapping
		pr.Approvers = make([]Identity, 0, len(m.Reviewers))
		for _, reviewer := range m.Reviewers {
			if isDirectlyApproved(reviewer) {
				pr.Approvers = append(pr.Approvers, toIdentity(&reviewer.IdentityRef))
			}
		}
		sort.Slice(pr.Approvers, func(i, j int) bool {
			return pr.Approvers[i].Name < pr.Approvers[j].Name
		})
		pr.CreatedBy = toIdentity(m.CreatedBy)

		// Add build status if available
		if evaluations != nil {
			if evals, ok := evaluations[m.PullRequestId]; ok {
				// Use repository from PR if available, otherwise use passed repo
				prRepo := m.Repository
				if prRepo == nil {
					prRepo = repo
				}
				pr.BuildStatus = parseBuildStatus(evals, orgName, prRepo)
				pr.PolicyChecks = resolvePolicyChecks(&m, evals)
			}
		}

		return pr
	}
}

// cleanBranchName removes the "refs/heads/" prefix from branch names
func cleanBranchName(refName string) string {
	return strings.TrimPrefix(refName, "refs/heads/")
}

// parseBuildStatus extracts the build validation status from policy evaluations.
// It looks for build validation policy evaluations (type ID: 0609b952-1397-4640-95ec-e00a01b2c241).
func parseBuildStatus(
	evaluations []models.PolicyEvaluationRecord,
	orgName string,
	repo *models.GitRepository,
) *BuildStatus {
	for _, eval := range evaluations {
		if eval.Configuration.Type.Id != models.PolicyTypeBuildValidation {
			continue
		}

		bs := &BuildStatus{TargetURL: extractBuildURL(eval, orgName, repo)}
		if displayName, ok := eval.Configuration.Settings["displayName"].(string); ok {
			bs.Description = displayName
		}
		bs.State = determineBuildStatus(eval)
		bs.Icon, bs.StatusText = getStatusDisplay(bs.State)
		return bs
	}
	return nil
}

// extractBuildURL constructs the build URL from the evaluation context
func extractBuildURL(
	eval models.PolicyEvaluationRecord,
	orgName string,
	repo *models.GitRepository,
) string {
	if eval.Context == nil || orgName == "" || repo == nil || repo.Project == nil {
		return ""
	}

	buildId := extractBuildId(eval.Context)
	if buildId <= 0 || repo.Project.Name == "" {
		return ""
	}

	return fmt.Sprintf("https://dev.azure.com/%s/%s/_build/results?buildId=%d",
		orgName, repo.Project.Name, buildId)
}

// extractBuildId attempts to extract the build ID from various possible locations in the context
func extractBuildId(context map[string]any) int {
	if id := convertToInt(context["buildId"]); id > 0 {
		return id
	}
	if build, ok := context["build"].(map[string]any); ok {
		if id := convertToInt(build["id"]); id > 0 {
			return id
		}
	}
	return convertToInt(context["id"])
}

// convertToInt converts various types to int
func convertToInt(v any) int {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return 0
}

// determineBuildStatus maps PolicyEvaluationStatus to PRBuildStatus
// TODO (tai): need to check again with expired build status
func determineBuildStatus(eval models.PolicyEvaluationRecord) PRBuildStatus {
	if eval.Configuration.Type.Id != models.PolicyTypeBuildValidation {
		return PRBuildStatusUnknown
	}

	switch eval.Status {
	case models.PolicyEvaluationStatusApproved:
		return PRBuildStatusSucceeded

	case models.PolicyEvaluationStatusRejected:
		return PRBuildStatusFailed

	case models.PolicyEvaluationStatusQueued:
		return determineQueuedStatus(eval)

	case models.PolicyEvaluationStatusRunning:
		return PRBuildStatusRunning

	case models.PolicyEvaluationStatusBroken:
		return PRBuildStatusError

	default:
		return PRBuildStatusUnknown
	}
}

// determineQueuedStatus distinguishes between expired and pending builds
func determineQueuedStatus(eval models.PolicyEvaluationRecord) PRBuildStatus {
	// An expired build has completed in the past but needs re-running
	hasCompletedDate := eval.CompletedDate != nil
	hasBuildId := extractBuildId(eval.Context) > 0

	if hasCompletedDate && hasBuildId {
		return PRBuildStatusExpired
	}
	return PRBuildStatusPending
}

// getStatusDisplay returns the icon and text for a given build status
func getStatusDisplay(status PRBuildStatus) (icon, text string) {
	switch status {
	case PRBuildStatusSucceeded:
		return styles.Success(ui.IconSuccess), styles.Success("passes")
	case PRBuildStatusFailed:
		return styles.Error(ui.IconFailure), styles.Error("fails")
	case PRBuildStatusRunning:
		return styles.Pending(ui.IconRunning), styles.Pending("running")
	case PRBuildStatusPending:
		return styles.Pending(ui.IconPending), styles.Pending("pending")
	case PRBuildStatusExpired:
		return styles.Warn(ui.IconWarning), styles.Warn("expired")
	case PRBuildStatusError:
		return styles.Warn(ui.IconWarning), styles.Warn("error")
	default:
		return "?", string(status)
	}
}

func editPrInfo(info *prInfo, editorCmd string) (*prInfo, error) {
	content := fmt.Sprintf("%s\n\n%s", info.title, info.desc)

	// Use the configured editor from global config, which handles fallbacks properly
	ed := editor.New("PR_EDIT*.md", editorCmd)

	updatedContent, err := ed.Edit(content)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(updatedContent, "\n\n", 2)
	newTitle := strings.TrimSpace(parts[0])
	if newTitle == "" {
		return nil, ErrEmptyTitle
	}
	newDesc := ""
	if len(parts) > 1 {
		newDesc = strings.TrimSpace(parts[1])
	}

	return &prInfo{title: newTitle, desc: newDesc}, nil
}

type prInfo struct {
	title string
	desc  string
}

// resolvePolicyChecks converts policy evaluations to simplified PolicyCheck structs
func resolvePolicyChecks(
	pr *models.GitPullRequest,
	evaluations []models.PolicyEvaluationRecord,
) PolicyChecks {
	if len(evaluations) == 0 {
		return nil
	}

	// Pre-allocate with estimated capacity
	checks := make([]PolicyCheck, 0, len(evaluations)+1)

	// Add a merge conflict check first if needed
	switch pr.MergeStatus {
	case models.PullRequestAsyncStatusConflicts:
		checks = append(checks, PolicyCheck{
			Name:       "Merge conflicts",
			Status:     "rejected",
			IsRequired: true,
			Icon:       styles.Error(ui.IconFailure),
		})
	case models.PullRequestAsyncStatusSucceeded:
		checks = append(checks, PolicyCheck{
			Name:       "No merge conflicts",
			Status:     "approved",
			IsRequired: true,
			Icon:       styles.Success(ui.IconSuccess),
		})
	}

	// Deduplicate by name and status - pre-allocate map
	type dedupeKey struct {
		name   string
		status string
	}
	seen := make(map[dedupeKey]bool, len(evaluations))

	for _, eval := range evaluations {
		name := getPolicyDisplayName(eval)
		status := policyStatusToString(eval.Status)
		key := dedupeKey{name, status}

		if !seen[key] {
			seen[key] = true
			checks = append(checks, PolicyCheck{
				Name:       name,
				Status:     status,
				IsRequired: eval.Configuration.IsBlocking,
				Icon:       getPolicyStatusIcon(eval.Status),
			})
		}
	}

	// Sort: required first, then by name
	sort.Slice(checks, func(i, j int) bool {
		if checks[i].IsRequired != checks[j].IsRequired {
			return checks[i].IsRequired
		}
		return checks[i].Name < checks[j].Name
	})

	return checks
}

// getPolicyDisplayName extracts the display name for a policy evaluation
func getPolicyDisplayName(eval models.PolicyEvaluationRecord) string {
	// Build Validation policies: use buildDefinitionName from context
	if eval.Configuration.Type.Id == models.PolicyTypeBuildValidation {
		if name, ok := eval.Context["buildDefinitionName"].(string); ok && name != "" {
			return name
		}
	}

	// Status policies: use defaultDisplayName or displayName
	if eval.Configuration.Type.Id == models.PolicyTypeStatus {
		if defaultName, ok := eval.Configuration.Settings["defaultDisplayName"].(string); ok &&
			defaultName != "" {
			return defaultName
		}
	}

	// Try various fallbacks
	if displayName, ok := eval.Configuration.Settings["displayName"].(string); ok &&
		displayName != "" {
		return displayName
	}
	if eval.Configuration.Type.DisplayName != "" {
		return eval.Configuration.Type.DisplayName
	}
	return "Unknown Policy"
}

// getPolicyStatusIcon returns a styled icon for a policy evaluation status
func getPolicyStatusIcon(status models.PolicyEvaluationStatus) string {
	switch status {
	case models.PolicyEvaluationStatusApproved:
		return styles.Success(ui.IconSuccess)
	case models.PolicyEvaluationStatusRejected:
		return styles.Error(ui.IconFailure)
	case models.PolicyEvaluationStatusRunning:
		return styles.Pending(ui.IconRunning)
	case models.PolicyEvaluationStatusQueued:
		return styles.Pending(ui.IconPending)
	case models.PolicyEvaluationStatusBroken:
		return styles.Warn(ui.IconWarning)
	case models.PolicyEvaluationStatusNotApplicable:
		return styles.Faint("-")
	default:
		return "?"
	}
}

// policyStatusToString converts PolicyEvaluationStatus to a simple string
func policyStatusToString(status models.PolicyEvaluationStatus) string {
	switch status {
	case models.PolicyEvaluationStatusApproved:
		return "approved"
	case models.PolicyEvaluationStatusRejected:
		return "rejected"
	case models.PolicyEvaluationStatusRunning:
		return "running"
	case models.PolicyEvaluationStatusQueued:
		return "queued"
	case models.PolicyEvaluationStatusBroken:
		return "broken"
	case models.PolicyEvaluationStatusNotApplicable:
		return "not_applicable"
	default:
		return "unknown"
	}
}
