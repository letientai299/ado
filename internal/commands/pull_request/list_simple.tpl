{{- /* gotype: []github.com/letientai299/ado/internal/commands/pull_request.PR */ -}}
{{- range $pr := .}}
- {{ if $pr.IsDraft }}{{warn "DRAFT"}} | {{end}}{{ $pr.Title | h1 }}
  {{ $pr.WebURL }}
  create by {{ $pr.CreatedBy.Name | person }} on {{ $pr.CreationDate | time}}
  {{- if gt (len $pr.Approvers) 0}}, approved by
    {{ range $i, $a := $pr.Approvers -}}
      - {{ $a.Name | person }}
    {{ end -}}
  {{end -}}
{{end -}}
