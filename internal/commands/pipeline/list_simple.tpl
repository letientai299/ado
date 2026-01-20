{{- /* gotype: []github.com/letientai299/ado/internal/commands/pipeline.Pipeline */ -}}
{{- range $p := .}}
- {{ $p.Name | h1 }}
  {{- if $p.YamlFilename }}
  YAML: {{ $p.YamlFilename }}
  {{- end }}
  {{ $p.WebURL }}
  {{- if eq $p.QueueStatus "disabled" }}
  Status: {{ warn "disabled" }}
  {{- end }}
{{end -}}
