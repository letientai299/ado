{{- /* gotype: github.com/letientai299/ado/internal/commands/pull_request.AnalysisResult */ -}}
{{ "PR Analysis" | h1 }}
{{ "Period" | heading }}: {{ .From | faint }} → {{ .To | faint }}

{{ "Overview" | heading }}
  Total     : {{ .Total }}
  Completed : {{ .Completed }}
  Active    : {{ .Active }}
  Abandoned : {{ .Abandoned }}
  Drafts    : {{ .DraftCount }} ({{ .DraftRatio }})
{{- if gt .Completed 0 }}

{{ "Active time (completed PRs)" | heading }}
  Average : {{ .AvgActiveTime }}
  Median  : {{ .MedianActiveTime }}
{{- end }}
{{- if gt (len .TopContributors) 0 }}

{{ "Top contributors" | heading }}
{{- range $i, $c := .TopContributors }}
  {{ printf "%2d." (add1 $i) }} {{ $c.Name | person }} — {{ $c.Count }} PR(s)
{{- end }}
{{- end }}
{{- if gt (len .TopReviewers) 0 }}

{{ "Top reviewers" | heading }}
{{- range $i, $r := .TopReviewers }}
  {{ printf "%2d." (add1 $i) }} {{ $r.Name | person }} — {{ $r.Count }} review(s)
{{- end }}
{{- end }}
