package pull_request

import (
	_ "embed"
	"fmt"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/letientai299/ado/internal/util/sh"
	"github.com/spf13/cobra"
)

//go:embed create.md
var createDoc string

const (
	defaultPrTitleTemplate = `{{.BranchName | replaceAll "/" "-"}}`
	defaultPrDescTemplate  = `{{range .Commits}}- {{.Subject}}
{{end}}` // the newline is crucial
)

type prTemplates struct {
	// Title template to generate PR title from commits.
	Title string `yaml:"title,omitempty" json:"title,omitempty"`
	// Desc template to generate PR description from commits.
	Desc string `yaml:"desc,omitempty" json:"desc,omitempty"`
}

type CreateConfig struct {
	Templates prTemplates `yaml:"templates,omitempty" json:"templates,omitempty"`
	Target    string      `yaml:"target,omitempty"    json:"target,omitempty"`

	yes     bool
	publish bool
	browse  bool
}

func (cc CreateConfig) OnResolved(c *cobra.Command) error {
	return nil
}

func createCmd() *cobra.Command {
	opts := defaultCreateConfig()

	c := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a pull request",
		Long:    createDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := newCommon(cmd, opts)
			if err != nil {
				return err
			}
			return newCreateProcessor(c).process()
		},
	}

	flags := c.Flags()
	flags.StringVarP(&opts.Target, "target", "t", opts.Target, "target branch")
	flags.BoolVarP(&opts.yes, "yes", "y", false, "skip all prompt and editor")
	flags.BoolVarP(&opts.publish, "publish", "p", false, "publish the PR")
	flags.BoolVarP(&opts.browse, "browse", "b", false, "open PR in browser after creating")

	return c
}

func defaultCreateConfig() *CreateConfig {
	opts := &CreateConfig{
		Templates: prTemplates{
			Title: defaultPrTitleTemplate,
			Desc:  defaultPrDescTemplate,
		},
		Target: "main",
	}

	config.Register(config.CommandConfig{
		Path:   "pull-request.create",
		Target: opts,
	})
	return opts
}

type createProcessor struct {
	*common[*CreateConfig]
}

func newCreateProcessor(c *common[*CreateConfig]) *createProcessor {
	return &createProcessor{common: c}
}

func (p *createProcessor) process() error {
	source, err := gitcli.CurrentBranch()
	if err != nil {
		return fmt.Errorf("fail to get current branch: %w", err)
	}

	if err = p.rebaseFromTarget(); err != nil {
		return err
	}

	if err = gitcli.SyncToRemote(source, p.confirmFn); err != nil {
		return err
	}

	existing, err := p.client.Git().PRs(p.cfg.Repository).List(p.ctx, git_prs.ListQuery{
		SearchCriteria: &git_prs.SearchCriteria{
			SourceRefName: util.Ptr("refs/heads/" + source),
			TargetRefName: util.Ptr("refs/heads/" + p.opts.Target),
			Status:        util.Ptr(models.PullRequestStatusActive),
		},
	})

	if err == nil && len(existing) > 0 {
		return p.updateExisting(existing[0])
	}

	return p.createNew(source, p.opts.Target)
}

func (p *createProcessor) rebaseFromTarget() error {
	target := p.opts.Target
	source, err := gitcli.CurrentBranch()
	if err != nil {
		return err
	}

	div, err := gitcli.CompareRevision(target, source)
	if err != nil {
		return err
	}

	if len(div.Behind) == 0 {
		return nil
	}

	ask := fmt.Sprintf("Current branch is %d commit(s) behind %s. Rebase?", len(div.Behind), target)
	if !p.confirmFn(ask) {
		return nil
	}

	if err = gitcli.Rebase(target); err != nil {
		return util.StrErr("rebase failed with conflicts, please resolve manually")
	}

	return nil
}

func (p *createProcessor) createNew(source, target string) error {
	info, err := p.genPrInfo(target, source)
	if err != nil {
		return fmt.Errorf("fail to generate PR details: %w", err)
	}

	if !p.opts.yes {
		if err = info.editWith(p.cfg.Editor); err != nil {
			return err
		}
	}

	pr := models.GitPullRequest{
		SourceRefName: "refs/heads/" + source,
		TargetRefName: "refs/heads/" + target,
		Title:         info.title,
		Description:   info.desc,
		IsDraft:       !p.opts.publish,
	}

	created, err := p.client.Git().PRs(p.cfg.Repository).Create(p.ctx, pr)
	if err != nil {
		return fmt.Errorf("fail to create PR: %w", err)
	}

	return p.postProcess(created)
}

func (p *createProcessor) genPrInfo(target, source string) (*prInfo, error) {
	commits, err := commitsAhead(target, source)
	if err != nil {
		return nil, err
	}

	info := &prInfo{commits: commits, isNew: true}

	if len(commits) == 1 {
		info.title = commits[0].Subject
		info.desc = commits[0].Body
		return info, nil
	}

	data := struct {
		BranchName string
		Commits    []gitcli.Commit
	}{
		BranchName: source,
		Commits:    commits,
	}

	if info.title, err = styles.RenderS(p.opts.Templates.Title, data); err != nil {
		return nil, err
	}

	if info.desc, err = styles.RenderS(p.opts.Templates.Desc, data); err != nil {
		return nil, err
	}

	return info, nil
}

func commitsAhead(target string, source string) ([]gitcli.Commit, error) {
	div, err := gitcli.CompareRevision(target, source)
	if err != nil {
		return nil, fmt.Errorf("fail to get commits: %w", err)
	}

	if len(div.Ahead) == 0 {
		return nil, fmt.Errorf("no commits found between %s and %s", target, source)
	}

	return div.Ahead, nil
}

func (p *createProcessor) confirmFn(ask string) bool {
	if p.opts.yes {
		return true
	}
	return ui.Confirm(ask, true)
}

func (p *createProcessor) updateExisting(pr models.GitPullRequest) error {
	commits, err := commitsAhead(pr.TargetRefName, cleanBranchName(pr.SourceRefName))
	if err != nil {
		return err
	}

	info := &prInfo{
		commits: commits,
		title:   pr.Title,
		desc:    pr.Description,
	}

	if err = info.editWith(p.cfg.Editor); err != nil {
		return err
	}

	updated, err := p.client.Git().PRs(p.cfg.Repository).Update(p.ctx,
		pr.PullRequestId,
		rest.PrUpdateRequest{
			Title:       util.Ptr(info.title),
			IsDraft:     util.Ptr(pr.IsDraft),
			Description: util.Ptr(info.desc),
		})
	if err != nil {
		return fmt.Errorf("fail to update PR: %w", err)
	}

	return p.postProcess(updated)
}

func (p *createProcessor) postProcess(pr *models.GitPullRequest) error {
	fmt.Printf("PR #%d: %s\n", pr.PullRequestId, styles.H1(pr.Title))
	url := fmt.Sprintf("%s/%d", p.baseURL, pr.PullRequestId)
	fmt.Println(url)

	if p.opts.browse {
		return sh.Browse(url)
	}

	return nil
}
