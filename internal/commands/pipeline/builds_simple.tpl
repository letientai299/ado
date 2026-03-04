{{- range $b := .}}
- #{{ $b.Number }} {{- if eq $b.Result "succeeded" }} {{ $b.Result | success }}
  {{- else if eq $b.Result "failed" }} {{ $b.Result | error }}
  {{- else if eq $b.Result "canceled" }} {{ $b.Result | warn }}
  {{- else if eq $b.Status "inProgress" }} {{ $b.Status | highlight }}
  {{- else }} {{ $b.Status }}
  {{- end }}
  ID: {{ $b.Id }} | {{ $b.Branch }} | {{ $b.Reason }}{{ if $b.Duration }} | {{ $b.Duration }}{{ end }}
  {{- if $b.Commit }}
  {{ $b.Commit }}
  {{- end }}
{{- end}}
