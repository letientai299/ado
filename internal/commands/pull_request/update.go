package pull_request

import (
	_ "embed"
	"fmt"
	"io"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/spf13/cobra"
)

//go:embed update.md
var updateDoc string

type UpdateConfig struct {
	filterConfig

	currentBranch bool
	edit          bool
	execute       *util.EnumFlag[action]
}

type updateData struct {
	updated bool
	model   *models.GitPullRequest
}

func updateCmd() *cobra.Command {
	opts := &UpdateConfig{
		execute: util.NewEnumFlag(allActions...).Optional(),
	}
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "Update a pull request",
		Long:    updateDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.execute.Validate(); err != nil {
				return err
			}

			opts.keywords = args
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}

			return newUpdateProcessor(c).process(args)
		},
	}

	opts.RegisterFlags(cmd)

	flags := cmd.Flags()
	flags.BoolVarP(
		&opts.currentBranch,
		"current-branch",
		".",
		opts.currentBranch,
		"only PRs of the current branch",
	)
	flags.BoolVarP(&opts.edit, "edit", "e", opts.edit, "edit title and description")

	flags.VarP(opts.execute, "execute", "x", "execute an action or modify the PR status")
	opts.execute.RegisterCompletion(cmd, "execute")
	return cmd
}

type updateProcessor struct {
	*common[*UpdateConfig]
	vp *viewProcessor
}

func newUpdateProcessor(c *common[*UpdateConfig]) *updateProcessor {
	vp := newViewProcessor(copyCommon(c, func(b *common[*ViewConfig]) *common[*ViewConfig] {
		b.opts = &ViewConfig{filterConfig: c.opts.filterConfig}
		return b
	}))
	return &updateProcessor{common: c, vp: vp}
}

func (u *updateProcessor) process(args []string) error {
	pr, err := u.findPR(args)
	if err != nil || pr == nil {
		return err
	}

	if err = u.vp.renderOne(*pr); err != nil {
		return err
	}

	data, err := u.prepareUpdateData(pr)
	if err != nil {
		return err
	}

	if !data.updated {
		return u.inform("No update", pr)
	}

	return u.updateToADO(pr.PullRequestId, data.model)
}

func (u *updateProcessor) findPR(args []string) (*models.GitPullRequest, error) {
	if u.opts.currentBranch {
		return u.findByCurrentBranch()
	}

	id, err := u.vp.findPrID(args)
	if err != nil || id == 0 {
		return nil, err
	}

	return u.client.Git().PRs(u.cfg.Repository).ByID(u.ctx, id)
}

func (u *updateProcessor) findByCurrentBranch() (*models.GitPullRequest, error) {
	branch, err := gitcli.CurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	refName := "refs/heads/" + branch
	criteria := &git_prs.SearchCriteria{
		SourceRefName: util.Ptr(refName),
		// `pr update` allows move reactive an abandoned PR, so we need to search for all status.
		// This search contains source branch info, so, practically, there won't be too many PRs.
		Status: util.Ptr(models.PullRequestStatusAll),
	}
	list, err := u.client.Git().PRs(u.cfg.Repository).List(u.ctx, git_prs.ListQuery{
		Top:            util.Ptr(20),
		SearchCriteria: criteria,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("no active pull request found for branch '%s'", branch)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	pr, _ := pick(list)
	return &pr, nil
}

func (u *updateProcessor) prepareUpdateData(pr *models.GitPullRequest) (*updateData, error) {
	data := &updateData{
		updated: false,
		model:   &models.GitPullRequest{},
	}

	if err := u.updateInfo(pr, data); err != nil {
		return nil, err
	}

	if act, ok := u.pickAction(pr); ok {
		data.updated = data.updated || act.exec(pr, data.model)
	}

	return data, nil
}

func (u *updateProcessor) updateInfo(m *models.GitPullRequest, data *updateData) error {
	if !u.opts.edit && !ui.Confirm("Edit PR title and description?", false) {
		return nil
	}

	curInfo := &prInfo{title: m.Title, desc: m.Description}
	newInfo, err := editPrInfo(curInfo, u.cfg.Editor)
	if err != nil {
		return err
	}

	if newInfo.title != curInfo.title {
		data.model.Title = newInfo.title
		data.updated = true
	}

	if newInfo.desc != curInfo.desc {
		data.model.Description = newInfo.desc
		data.updated = true
	}

	return nil
}

func (u *updateProcessor) updateToADO(id int32, data *models.GitPullRequest) error {
	updated, err := u.client.Git().PRs(u.cfg.Repository).Update(u.ctx, id, *data)
	if err != nil {
		return err
	}
	return u.inform("Updated", updated)
}

func (u *updateProcessor) inform(msg string, pr *models.GitPullRequest) error {
	log.Infof("%s, #%d: %s", msg, pr.PullRequestId, styles.H1(pr.Title))
	_, err := fmt.Println(webURL(u.baseURL, pr.PullRequestId))
	return err
}

func (u *updateProcessor) pickAction(cur *models.GitPullRequest) (action, bool) {
	act := u.opts.execute.Value()
	if act != "" {
		return act, true
	}

	usable := slices.DeleteFunc(allActions, func(a action) bool { return !a.applicable(cur) })
	picked := ui.Pick[action](usable, ui.PickConfig[action]{
		Title:       "Which action? (ctrl-c to cancel)",
		Render:      func(w io.Writer, it action, _ []int) { _, _ = w.Write([]byte(it)) },
		FilterValue: func(item action) string { return string(item) },
	})

	if picked.IsNil() {
		return "", false
	}

	return picked.Get(), true
}
