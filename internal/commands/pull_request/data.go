package pull_request

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/letientai299/ado/internal/util/fp"
)

type Vote string

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

func isApproved(vote *models.IdentityRefWithVote) bool {
	// NOTE (tai): there's Approve (vote=10) and Approve-with-suggestions (vote=5).
	//  Almost no one use approve with suggestions, so we only consider vote > 0 as approval.
	//  We might update this logic if there's a need later
	return vote.Vote > 0
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

		approvers := fp.Map(
			slices.DeleteFunc(m.Reviewers, isApproved),
			func(x *models.IdentityRefWithVote) *models.IdentityRef { return &x.IdentityRef },
		)
		pr.Approvers = fp.Map(approvers, toIdentity)
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
			}
		}

		return pr
	}
}

// cleanBranchName removes the "refs/heads/" prefix from branch names
func cleanBranchName(refName string) string {
	const prefix = "refs/heads/"
	if strings.HasPrefix(refName, prefix) {
		return refName[len(prefix):]
	}
	return refName
}

// parseBuildStatus extracts the build validation status from policy evaluations.
// It looks for build validation policy evaluations (type ID: 0609b952-1397-4640-95ec-e00a01b2c241).
func parseBuildStatus(
	evaluations []models.PolicyEvaluationRecord,
	orgName string,
	repo *models.GitRepository,
) *BuildStatus {
	if len(evaluations) == 0 {
		return nil
	}

	// Look for build validation policy
	for _, eval := range evaluations {
		if eval.Configuration.Type.Id == models.PolicyTypeBuildValidation {
			return buildStatusFromEvaluation(eval, orgName, repo)
		}
	}

	// No build validation policy found
	return nil
}

// buildStatusFromEvaluation converts a PolicyEvaluationRecord to BuildStatus
func buildStatusFromEvaluation(
	eval models.PolicyEvaluationRecord,
	orgName string,
	repo *models.GitRepository,
) *BuildStatus {
	bs := &BuildStatus{
		Description: extractDisplayName(eval),
		TargetURL:   extractBuildURL(eval, orgName, repo),
	}

	// Determine the build status and set corresponding icon and text
	bs.State = determineBuildStatus(eval)
	bs.Icon, bs.StatusText = getStatusDisplay(bs.State)

	return bs
}

// extractDisplayName gets the display name from the evaluation settings
func extractDisplayName(eval models.PolicyEvaluationRecord) string {
	if displayName, ok := eval.Configuration.Settings["displayName"].(string); ok {
		return displayName
	}
	return ""
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
	// Try different possible locations for buildId
	candidates := []any{
		context["buildId"],
		extractFromBuildObject(context),
		context["id"],
	}

	for _, candidate := range candidates {
		if id := convertToInt(candidate); id > 0 {
			return id
		}
	}

	return 0
}

// extractFromBuildObject tries to extract ID from a nested build object
func extractFromBuildObject(context map[string]any) any {
	if build, ok := context["build"].(map[string]interface{}); ok {
		return build["id"]
	}
	return nil
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
func determineBuildStatus(eval models.PolicyEvaluationRecord) PRBuildStatus {
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
		return styles.Faint(ui.IconRunning), styles.Time("running")
	case PRBuildStatusPending:
		return styles.Faint(ui.IconPending), styles.Faint("pending")
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
