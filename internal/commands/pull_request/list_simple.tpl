{{- /* gotype: []github.com/letientai299/ado/internal/commands/pull_request.PR */ -}}
{{- range $pr := .}}
{{ if $pr.IsDraft }}{{warn "DRAFT"}} | {{end}}{{ $pr.Title | h1 }}
  {{ $pr.SourceBranchName | const }} --> {{ $pr.TargetBranchName | const }}
  {{- if $pr.BuildStatus }}
  {{ $pr.BuildStatus.Icon }} Build {{ $pr.BuildStatus.StatusText }}: {{ $pr.BuildStatus.TargetURL }}
  {{- end }}
  Created by {{ $pr.CreatedBy.Name | person }} on {{ $pr.CreationDate | time}}
  {{- if gt (len $pr.Approvers) 0}}
  Approved by {{ range $i, $a := $pr.Approvers -}}{{ $a.Name | person }}, {{ end -}}{{end}}
  PR Link: {{ $pr.WebURL }}
{{end -}}
