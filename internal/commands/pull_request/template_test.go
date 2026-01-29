package pull_request

import (
	"testing"

	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/stretchr/testify/require"
)

func TestTemplates(t *testing.T) {
	data := struct {
		BranchName string
		Commits    []gitcli.Commit
	}{
		BranchName: "feature/abc-123",
		Commits: []gitcli.Commit{
			{Subject: "feat: first commit"},
			{Subject: "fix: second commit"},
		},
	}

	t.Run("DefaultTitleTemplate", func(t *testing.T) {
		got, err := styles.RenderS(defaultPrTitleTemplate, data)
		require.NoError(t, err)
		require.Equal(t, "feature abc 123", got)
	})

	t.Run("DefaultDescTemplate", func(t *testing.T) {
		got, err := styles.RenderS(defaultPrDescTemplate, data)
		require.NoError(t, err)
		// Current default template: {{range .Commits}}- {{.Subject}}\n{{end}}
		want := "- feat: first commit\n- fix: second commit\n"
		require.Equal(t, want, got)
	})

	t.Run("CreateMdExampleTitleTemplate", func(t *testing.T) {
		// Example from create.md: {{.BranchName | replaceAll "/" "-"}}
		tpl := `{{.BranchName | replaceAll "/" "-"}}`
		got, err := styles.RenderS(tpl, data)
		require.NoError(t, err)
		require.Equal(t, "feature-abc-123", got)
	})

	t.Run("CreateMdExampleDescTemplate", func(t *testing.T) {
		// Example from create.md:
		// {{range .Commits}}- {{.Subject}}
		// {{end}}
		tpl := "{{range .Commits}}- {{.Subject}}\n{{end}}"
		got, err := styles.RenderS(tpl, data)
		require.NoError(t, err)
		want := "- feat: first commit\n- fix: second commit\n"
		require.Equal(t, want, got)
	})

	t.Run("CreateMdYamlExampleDescTemplate", func(t *testing.T) {
		// Example from create.md YAML block:
		tpl := `{{- range .Commits -}}
- {{ .Subject }}
{{- if .Body }}

{{ .Body | indent 2 }}
{{- end }}

{{ end -}}`
		dataWithBody := struct {
			BranchName string
			Commits    []gitcli.Commit
		}{
			BranchName: "feature/abc-123",
			Commits: []gitcli.Commit{
				{
					Subject: "commit 1 subject line",
					Body:    "Body 1 line 1, might be multi line",
				},
				{
					Subject: "commit 2 subject line",
					Body: `- list item 1

` + "```cs" + `
// some code
` + "```" + `

- list item 2
- list item 3
  - sub list item 3`,
				},
				{
					Subject: "commit 3 subject line",
				},
			},
		}
		got, err := styles.RenderS(tpl, dataWithBody)
		require.NoError(t, err)

		want := `- commit 1 subject line

  Body 1 line 1, might be multi line

- commit 2 subject line

  - list item 1

  ` + "```cs" + `
  // some code
  ` + "```" + `

  - list item 2
  - list item 3
    - sub list item 3

- commit 3 subject line

`
		require.Equal(t, want, got)
		t.Log(got)
	})
}
