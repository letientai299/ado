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
	PullRequestId int32      `yaml:"pull_request_id" json:"pull_request_id,omitempty"`
	Title         string     `yaml:"title"           json:"title,omitempty"`
	Description   string     `yaml:"description"     json:"description,omitempty"`
	IsDraft       bool       `yaml:"is_draft"        json:"is_draft,omitempty"`
	CreatedBy     Identity   `yaml:"created_by"      json:"created_by"`
	CreationDate  string     `yaml:"creation_date"   json:"creation_date,omitempty"`
	WebURL        string     `yaml:"web_url"         json:"web_url,omitempty"`
	Approvers     []Identity `yaml:"approvers"       json:"approvers,omitempty"`
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

func converter(baseURL string) func(m models.GitPullRequest) PR {
	return func(m models.GitPullRequest) PR {
		pr := PR{
			PullRequestId: m.PullRequestId,
			Title:         m.Title,
			Description:   m.Description,
			IsDraft:       m.IsDraft,
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
		return pr
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
