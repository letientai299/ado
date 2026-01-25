package pull_request

import (
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/letientai299/ado/internal/util/gitcli"
)

const (
	ErrEmptyTitle util.StrErr = "PR title cannot be empty."
	ErrUnchanged  util.StrErr = "PR title and description unchanged."
	blankLine                 = "\n\n"
)

const (
	// prEditingMarker is used to separate the PR title/desc from the commit messages reference
	prEditingMarker = "<!-- ado-pr-editing: DO NOT REMOVE -->"

	// prReferences is the template for displaying commit messages as reference when editing
	prReferences = `<!-- All commit messages for your references -->
<!--
{{ range .}}
## {{.Subject}}{{if .Body}}
{{.Body}}{{end}}
{{end -}}
-->
`
)

type prInfo struct {
	commits []gitcli.Commit
	title   string
	desc    string
	isNew   bool
}

func (p *prInfo) editWith(editorCmd string) error {
	ref, err := p.renderCommitRefs()
	if err != nil {
		return err
	}

	content := fmt.Sprintf("%s\n\n%s\n%s", p.title, p.desc, ref)
	updatedContent, err := editor.New("PR_EDIT*.md", editorCmd).Edit(content)
	if err != nil {
		return err
	}
	return p.parse(updatedContent)
}

func (p *prInfo) renderCommitRefs() (string, error) {
	if len(p.commits) == 0 {
		return "", nil
	}
	ref, err := styles.RenderS(prReferences, p.commits)
	if err != nil {
		return "", fmt.Errorf("failed to render diff template: %w", err)
	}

	return prEditingMarker + "\n" + ref, nil
}

func (p *prInfo) parse(content string) error {
	s, _, _ := strings.Cut(content, prEditingMarker) // drop everything after the marker
	s = strings.TrimSpace(s)

	title, desc, _ := strings.Cut(s, blankLine)
	if title == "" {
		return ErrEmptyTitle
	}

	if title == p.title && desc == p.desc && !p.isNew {
		return ErrUnchanged
	}

	p.title = title
	p.desc = desc
	return nil
}
