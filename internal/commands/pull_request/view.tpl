{{- .Title | h1 }}
{{.SourceBranchName | const}} --> {{.TargetBranchName | const}}
{{if .BuildStatus -}}
{{if eq .BuildStatus.StatusText "passes" }}{{ .BuildStatus.Icon | success }}{{ else }}{{ .BuildStatus.Icon | warn }}{{ end }} Build {{ .BuildStatus.StatusText }}: {{ .BuildStatus.TargetURL }}
{{end -}}
{{if gt (len .Description) 0 -}}
{{.Description | trimSpace | markdown -}}
{{- end}}
{{- "PR Link:" | heading}} {{.WebURL}}
