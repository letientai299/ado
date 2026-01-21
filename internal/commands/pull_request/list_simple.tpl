{{- /* gotype: []github.com/letientai299/ado/internal/commands/pull_request.PR */ -}}
{{- range $pr := .}}
- {{ if $pr.IsDraft }}{{warn "DRAFT"}} | {{end}}{{ $pr.Title | h1 }}
  {{ $pr.SourceBranchName | const }} --> {{ $pr.TargetBranchName | const }}
  {{- if $pr.BuildStatus }}
  {{ if eq $pr.BuildStatus.StatusText "passes" }}{{ $pr.BuildStatus.Icon | success }}{{ else }}{{ $pr.BuildStatus.Icon }}{{ end }} Build {{ $pr.BuildStatus.StatusText }}: {{ $pr.BuildStatus.TargetURL }}
  {{- end }}
  {{ $pr.WebURL }}
  create by {{ $pr.CreatedBy.Name | person }} on {{ $pr.CreationDate | time}}
  {{- if gt (len $pr.Approvers) 0}}, approved by
    {{ range $i, $a := $pr.Approvers -}}
      - {{ $a.Name | person }}
    {{ end -}}
  {{end -}}
{{end -}}
