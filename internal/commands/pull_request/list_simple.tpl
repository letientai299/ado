{{- /* gotype: []github.com/letientai299/ado/internal/commands/pull_request.PR */ -}}
{{- range $pr := .}}
- {{ if $pr.IsDraft }}{{warn "DRAFT"}} | {{end}}{{ $pr.Title | heading }}
  create by {{ $pr.CreatedBy.DisplayName | person }} at {{ $pr.CreationDate | time}}
  {{ $pr.WebURL }}
{{end -}}
