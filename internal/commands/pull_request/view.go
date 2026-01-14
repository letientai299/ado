package pull_request

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

var prView = &cobra.Command{
	Use:     "view",
	Aliases: []string{"view", "v"},
	Short:   "View detail of a pull request",
	RunE: func(cmd *cobra.Command, args []string) error {
		return View(cmd.Context())
	},
}

func View(ctx context.Context) error {
	cfg := config.From(ctx)
	pr, err := rest.New(cfg.Tenant).
		Git().
		PRs(cfg.Repository).
		ByID(ctx, 1329796)
	if err != nil {
		log.Error(err)
		return err
	}

	return util.DumpJSON(pr)
}
