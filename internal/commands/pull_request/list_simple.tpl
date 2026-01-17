{{- /* gotype: []github.com/letientai299/ado/internal/models.GitPullRequest */ -}}
{{- range $pr := .}}
- {{ if $pr.IsDraft }}{{warn "DRAFT"}} | {{end}}{{ $pr.Title | heading }}
  create by {{ $pr.CreatedBy.DisplayName | person }} at {{ $pr.CreationDate.Format "2006-01-02" | time}}
  {{ $pr | webURL -}}
{{end -}}
