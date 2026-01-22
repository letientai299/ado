package pull_request

import (
	"strings"
	"testing"

	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/stretchr/testify/assert"
)

func TestGenPrInfo(t *testing.T) {
	defaultOpts := &CreateConfig{
		Templates: prTemplates{
			Title: defaultPrTitleTemplate,
			Desc:  defaultPrDescTemplate,
		},
	}
	tests := []struct {
		name    string
		branch  string
		commits []gitcli.Commit
		opts    *CreateConfig
		want    *prInfo
		wantErr bool
	}{
		{
			name:   "single commit",
			branch: "feature/foo",
			commits: []gitcli.Commit{
				{Subject: "feat: add foo", Body: "details about foo"},
			},
			opts: defaultOpts,
			want: &prInfo{
				title: "feat: add foo",
				desc:  "details about foo",
			},
			wantErr: false,
		},
		{
			name:   "multiple commits - default templates",
			branch: "feature/foo-bar",
			commits: []gitcli.Commit{
				{Subject: "feat: add foo", Body: "details about foo"},
				{Subject: "fix: bar bug", Body: "details about bar"},
			},
			opts: defaultOpts,
			want: &prInfo{
				title: "feature-foo-bar",
				desc:  "- feat: add foo\n- fix: bar bug",
			},
			wantErr: false,
		},
		{
			name:   "multiple commits - custom templates",
			branch: "feature/custom",
			commits: []gitcli.Commit{
				{Subject: "feat: 1", Body: "b1"},
				{Subject: "feat: 2", Body: "b2"},
			},
			opts: &CreateConfig{
				Templates: prTemplates{
					Title: "PR: {{.BranchName}}",
					Desc:  "Commits: {{len .Commits}}",
				},
			},
			want: &prInfo{
				title: "PR: feature/custom",
				desc:  "Commits: 2",
			},
			wantErr: false,
		},
		{
			name:   "multiple commits - custom markdown templates from create.md",
			branch: "prefix/task-1",
			commits: []gitcli.Commit{
				{Subject: "feat: add foo", Body: "details about foo"},
				{Subject: "fix: bar bug", Body: ""},
			},
			opts: &CreateConfig{
				Templates: prTemplates{
					Title: `{{ .BranchName | trimPrefix "prefix/"  |replaceAll "/" "-" }}`,
					Desc: strings.TrimSpace(`
{{range .Commits}}
- {{.Subject}}

  {{if .Body}}{{.Body}}{{end}}
{{end}}`),
				},
			},
			want: &prInfo{
				title: "task-1",
				desc: strings.TrimSpace(`
- feat: add foo

  details about foo

- fix: bar bug`),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &createProcessor{
				common: &common[*CreateConfig]{
					opts: tt.opts,
				},
			}
			got, err := p.genPrInfo(tt.branch, tt.commits)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
