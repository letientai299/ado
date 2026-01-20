package pull_request

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/sh"
	"github.com/spf13/cobra"
)

//go:embed view.tpl
var viewTpl string

//go:embed view.md
var viewDoc string

type ViewConfig struct {
	filterConfig
	browse bool
}

func viewCmd() *cobra.Command {
	opts := &ViewConfig{}

	cmd := &cobra.Command{
		Use:     "view <id|text>",
		Aliases: []string{"v"},
		Short:   "View detail of a pull request",
		Long:    viewDoc,
		Args:    cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}

			return newViewProcessor(c).process(args)
		},
	}
	opts.RegisterFlags(cmd)
	cmd.Flags().BoolVarP(&opts.browse, "browse", "b", false, "open PR in browser")
	return cmd
}

func newViewProcessor(c *common[*ViewConfig]) *viewProcessor {
	lp := newListProcessor(copyCommon(c, func(b *common[*ListConfig]) *common[*ListConfig] {
		b.opts = &ListConfig{filterConfig: c.opts.filterConfig}
		return b
	}))
	return &viewProcessor{common: c, lp: lp}
}

type viewProcessor struct {
	*common[*ViewConfig]
	lp listProcessor
}

func (v viewProcessor) process(args []string) error {
	prId, err := v.findPrID(args)
	if err != nil || prId == 0 {
		return err
	}

	return v.renderByID(prId)
}

func (v viewProcessor) findPrID(args []string) (int32, error) {
	// 1. Try if the first arg is a PR ID
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			var m *models.GitPullRequest
			// TODO (tai): in case of valid ID, we call ADO twice, should add ctx-cache,
			//  but be careful to not serving stale data in long running TUI
			m, err = v.client.Git().PRs(v.cfg.Repository).ByID(v.ctx, int32(id))
			if err == nil {
				return m.PullRequestId, nil
			}
			// if error, treat the numeric arg as a keyword
		}
	}

	// 2. Fallback to list/filter logic
	prs, err := v.lp.find()
	if err != nil {
		return 0, err
	}
	if len(prs) == 0 {
		return 0, errors.New("no pull request found matching the criteria")
	}

	if len(prs) == 1 {
		return prs[0].PullRequestId, nil
	}

	if pr, ok := pick(prs); ok {
		return pr.PullRequestId, nil
	}

	return 0, nil
}

const prPickTpl = `{{.Title}} ({{.CreatedBy.DisplayName|person}}, {{.CreationDate.Format "2016-01-16" | time }}{{if .IsDraft}}, {{warn "DRAFT"}}{{end}})`

func pick(prs []models.GitPullRequest) (models.GitPullRequest, bool) {
	selected := ui.Pick(prs, ui.PickConfig[models.GitPullRequest]{
		Render: func(w io.Writer, pr models.GitPullRequest, matches []int) {
			pr.Title = styles.HighlightMatch(pr.Title, matches)
			util.PanicIf(styles.Render(w, prPickTpl, pr))
		},
		FilterValue: func(pr models.GitPullRequest) string { return strings.ToLower(pr.Title) },
	})

	if selected.IsNil() {
		return models.GitPullRequest{}, false
	}

	return selected.Get(), true
}

func (v viewProcessor) renderByID(id int32) error {
	// use this ByID API to fetch full PR details.
	// The List API returns only max 400 chars for PR description.
	m, err := v.client.Git().PRs(v.cfg.Repository).ByID(v.ctx, id)
	if err != nil {
		return err
	}

	return v.renderOne(*m)
}

func (v viewProcessor) renderOne(m models.GitPullRequest) error {
	pr := converter(v.baseURL)(m)
	if v.opts.browse {
		fmt.Println(pr.WebURL)
		return sh.Browse(pr.WebURL)
	}
	return styles.RenderOut(viewTpl, pr)
}
