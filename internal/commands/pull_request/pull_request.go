package pull_request

import (
	_ "embed"

	"github.com/letientai299/ado/internal/models"
	"github.com/spf13/cobra"
)

//go:embed pull_request.md
var doc string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull-request",
		Aliases: []string{"pr", "pull"},
		Short:   "List, view, create or manipulate pull requests",
		Long:    doc,
	}
	cmd.AddCommand(
		listCmd(),
		viewCmd(),
		createCmd(),
		updateCmd(),
	)
	return cmd
}

type Vote string

type Identity struct {
	Id    string
	Name  string
	Email string
}

type PR struct {
	PullRequestId int        `yaml:"pull_request_id" json:"pull_request_id,omitempty"`
	Title         string     `yaml:"title"           json:"title,omitempty"`
	Description   string     `yaml:"description"     json:"description,omitempty"`
	IsDraft       bool       `yaml:"is_draft"        json:"is_draft,omitempty"`
	CreatedBy     Identity   `yaml:"created_by"      json:"created_by"`
	CreationDate  string     `yaml:"creation_date"   json:"creation_date,omitempty"`
	WebURL        string     `yaml:"web_url"         json:"web_url,omitempty"`
	Approvers     []Identity `yaml:"approvers"       json:"approvers,omitempty"`
}

func isApproved(vote models.IdentityRefWithVote) bool {
	// NOTE (tai): there's Approve (vote=10) and Approve-with-suggestions (vote=5).
	//  Almost no one use approve with suggestions, so we only consider vote > 0 as approval.
	//  We might update this logic if there's a need later
	return vote.Vote > 0
}

func toIdentity(a models.IdentityRef) Identity {
	return Identity{
		Id:    a.Id,
		Name:  a.DisplayName,
		Email: a.UniqueName,
	}
}
