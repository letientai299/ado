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

type PR struct {
	PullRequestId int
	Title         string
	Description   string
	IsDraft       bool
	CreatedBy     models.IdentityRef
	CreationDate  string
	WebURL        string
}
