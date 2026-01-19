package pull_request

import (
	_ "embed"
	"errors"
	"fmt"
	"strconv"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/styles"
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
		Args:    cobra.MinimumNArgs(1),
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
	lp := listProcessor{
		common: &common[*ListConfig]{
			ctx:     c.ctx,
			cfg:     c.cfg,
			client:  c.client,
			baseURL: c.baseURL,
			opts: &ListConfig{
				filterConfig: c.opts.filterConfig,
			},
		},
	}

	return &viewProcessor{common: c, lp: lp}
}

type viewProcessor struct {
	*common[*ViewConfig]
	lp listProcessor
}

func (v viewProcessor) process(args []string) error {
	// 1. Try if the first arg is a PR ID
	if len(args) == 1 {
		if id, err := strconv.ParseInt(args[0], 10, 32); err == nil {
			var m *models.GitPullRequest
			m, err = v.client.Git().PRs(v.cfg.Repository).ByID(v.ctx, int32(id))
			if err == nil {
				return v.renderOne(*m)
			}
		}
	}

	// 2. Fallback to list/filter logic
	prs, err := v.lp.find()
	if err != nil {
		return err
	}

	switch len(prs) {
	case 0:
		return errors.New("no pull request found matching the criteria")
	case 1:
		return v.renderByID(prs[0].PullRequestId)
	default:
		for _, pr := range prs {
			fmt.Printf("%s\t%s\n", pr.Title, pr.WebURL)
		}
		return nil
	}
}

func (v viewProcessor) renderByID(id int32) error {
	m, err := v.client.Git().PRs(v.cfg.Repository).ByID(v.ctx, id)
	if err != nil {
		return err
	}

	return v.renderOne(*m)
}

func (v viewProcessor) renderOne(m models.GitPullRequest) error {
	pr := v.lp.toPR(m)
	if v.opts.browse {
		fmt.Println(pr.WebURL)
		return sh.Browse(pr.WebURL)
	}
	return styles.RenderOut(viewTpl, pr)
}
