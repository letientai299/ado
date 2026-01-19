package pull_request

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/git_prs"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/letientai299/ado/internal/util/sh"
	"github.com/spf13/cobra"
)

//go:embed create.md
var createDoc string

const (
	defaultPrTitleTemplate = `{{replaceAll .BranchName "/" "-"}}`
	defaultPrDescTemplate  = `{{range .Commits}}- {{.Subject}}
{{end}}`
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
		return p.updateExistingPrInfo(existing[0])
	}

	return p.createNew(source, p.opts.Target)
}

func (p *createProcessor) createNew(source, target string) error {
	div, err := gitcli.CompareRevision(target, source)
	if err != nil {
		return fmt.Errorf("fail to get commits: %w", err)
	}

	info, err := p.genPrInfo(source, div.Ahead)
	if err != nil {
		return fmt.Errorf("fail to generate PR details: %w", err)
	}

	if !p.opts.yes {
		info, err = p.editPrInfo(info)
		if err != nil {
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

type prInfo struct {
	title string
	desc  string
}

func (p *createProcessor) genPrInfo(branch string, commits []gitcli.Commit) (*prInfo, error) {
	if len(commits) == 1 {
		cm := commits[0]
		return &prInfo{title: cm.Subject, desc: cm.Body}, nil
	}

	data := struct {
		BranchName string
		Commits    []gitcli.Commit
	}{
		BranchName: branch,
		Commits:    commits,
	}

	var title, desc string
	var err error
	if title, err = styles.RenderS(p.opts.Templates.Title, data); err != nil {
		return nil, err
	}

	if desc, err = styles.RenderS(p.opts.Templates.Desc, data); err != nil {
		return nil, err
	}

	return &prInfo{
		title: strings.TrimSpace(title),
		desc:  strings.TrimSpace(desc),
	}, nil
}

func (p *createProcessor) editPrInfo(info *prInfo) (*prInfo, error) {
	content := fmt.Sprintf("%s\n\n%s", info.title, info.desc)

	// Use the configured editor from global config, which handles fallbacks properly
	ed := editor.New("PR_EDIT*.md", p.cfg.Editor)

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

func (p *createProcessor) confirmFn(ask string) bool {
	if p.opts.yes {
		return true
	}
	return ui.Confirm(ask, true)
}

func (p *createProcessor) updateExistingPrInfo(pr models.GitPullRequest) error {
	var err error
	info := &prInfo{title: pr.Title, desc: pr.Description}

	info, err = p.editPrInfo(info)
	if err != nil {
		return fmt.Errorf("failed to edit PR info: %w", err)
	}

	updated, err := p.client.Git().PRs(p.cfg.Repository).Update(p.ctx,
		pr.PullRequestId,
		models.GitPullRequest{
			Title:       info.title,
			Description: info.desc,
		})
	if err != nil {
		return fmt.Errorf("fail to update PR: %w", err)
	}

	return p.postProcess(updated)
}

func (p *createProcessor) postProcess(pr *models.GitPullRequest) error {
	fmt.Printf("PR #%d: %s\n", pr.PullRequestId, styles.H1(pr.Title))
	webURL := fmt.Sprintf("%s/%d", p.baseURL, pr.PullRequestId)
	fmt.Println(webURL)

	if p.opts.browse {
		return sh.Browse(webURL)
	}

	return nil
}
