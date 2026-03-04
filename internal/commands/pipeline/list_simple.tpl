{{- range $p := .}}
- {{ $p.Name | h1 }} {{- if or (eq $p.QueueStatus "disabled") (eq $p.QueueStatus "paused") }} ({{printf "%s" $p.QueueStatus|warn}}){{- end }}
  ID: {{ $p.Id }}
  {{- if $p.YamlFilename }}
  YAML: {{ $p.YamlFilename }}
  {{- end }}
  {{ $p.WebURL }}
{{- end}}
