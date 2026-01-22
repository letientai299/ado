{{- .Title | h1 }}
{{ .SourceBranchName | const}} --> {{.TargetBranchName | const}}
{{if .BuildStatus -}}
{{.BuildStatus.Icon }} Build {{ .BuildStatus.StatusText }}: {{ .BuildStatus.TargetURL }}
{{- end -}}
{{- if gt (len .Description) 0 -}}
{{.Description | trimSpace | markdown -}}
{{- end -}}
{{- "PR Link:" | heading}} {{.WebURL}}
