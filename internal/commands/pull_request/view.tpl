{{- .Title | h1 }}
{{ .SourceBranchName }} --> {{.TargetBranchName }}
{{- if .BuildStatus }}
{{ .BuildStatus.Icon }} Build {{ .BuildStatus.StatusText }}: {{ .BuildStatus.TargetURL }}
{{- end }}
{{- if .PolicyChecks.Pending }}

{{"Pending policies" | heading}}
{{- range .PolicyChecks.Pending }}
  {{.Icon}} {{.Name}}
{{- end }}
{{- end }}
{{- if .Description }}

{{.Description | trimSpace | markdown -}}
{{- end }}
{{"PR Link:" | heading}} {{.WebURL}}
