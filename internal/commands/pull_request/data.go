package pull_request

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
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

// BuildStatus represents the build/pipeline status for a pull request
type BuildStatus struct {
	State       string `yaml:"state"       json:"state,omitempty"`       // succeeded, failed, pending, error
	Description string `yaml:"description" json:"description,omitempty"` // human-readable description
	TargetURL   string `yaml:"target_url"  json:"target_url,omitempty"`  // link to build results
	Icon        string `yaml:"icon"        json:"icon,omitempty"`        // emoji icon
	StatusText  string `yaml:"status_text" json:"status_text,omitempty"` // "passes", "fails", "pending"
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
		State: eval.Status,
	}

	// Get display name from settings
	if displayName, ok := eval.Configuration.Settings["displayName"].(string); ok {
		bs.Description = displayName
	}

	// Extract buildId from context and construct URL
	if eval.Context != nil {
		// Try to extract buildId
		var buildId interface{}
		if id, ok := eval.Context["buildId"]; ok {
			buildId = id
		} else if build, ok := eval.Context["build"].(map[string]interface{}); ok {
			buildId = build["id"]
		} else if id, ok := eval.Context["id"]; ok {
			buildId = id
		}

		// Convert buildId to int and construct URL
		if buildId != nil && orgName != "" && repo != nil && repo.Project != nil {
			var buildIdNum int
			switch v := buildId.(type) {
			case float64:
				buildIdNum = int(v)
			case int:
				buildIdNum = v
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					buildIdNum = parsed
				}
			}

			if buildIdNum > 0 && repo.Project.Name != "" {
				bs.TargetURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_build/results?buildId=%d",
					orgName, repo.Project.Name, buildIdNum)
			}
		}
	}

	// Set icon and status text based on status
	switch strings.ToLower(eval.Status) {
	case "approved":
		bs.Icon = "✓"
		bs.StatusText = "passes"
	case "rejected":
		bs.Icon = "✗"
		bs.StatusText = "fails"
	case "queued", "running":
		bs.Icon = "⏳"
		bs.StatusText = "pending"
	case "broken":
		bs.Icon = "⚠"
		bs.StatusText = "error"
	default:
		bs.Icon = "?"
		bs.StatusText = eval.Status
	}

	return bs
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
