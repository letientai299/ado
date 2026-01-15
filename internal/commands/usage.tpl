{{- /*gotype:github.com/spf13/cobra.Command */ -}}
{{- if gt (len .Aliases) 1 -}}
{{- "Aliases:" | headingStyle }} {{ range $i, $alias := .Aliases }}{{ if $i }}, {{ end }}{{ $alias | cmdStyle }}{{ end }}
{{- end }}

{{- if .HasExample }}

{{ "Examples:" | headingStyle }}
{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
  {{- $cmds := .Commands }}
  {{- if eq (len .Groups) 0 }}

{{ "Available Commands:" | headingStyle }}
    {{- range $cmds }}
      {{- if (or .IsAvailableCommand (eq .Name "help")) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
      {{- end }}
    {{- end }}
  {{- else }}
    {{- range $group := .Groups }}

{{ .Title | headingStyle }}
      {{- range $cmds }}
        {{- if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help"))) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
        {{- end }}
      {{- end }}
    {{- end }}

    {{- if not .AllChildCommandsHaveGroup }}

{{ "Additional Commands:" | headingStyle }}
      {{- range $cmds }}
        {{- if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help"))) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end -}}

{{- renderFlags . -}}

{{- if .HasHelpSubCommands }}

{{ "Additional help topics:" | headingStyle }}
  {{- range .Commands }}
    {{- if .IsAdditionalHelpTopicCommand }}
  {{ rpad .CommandPath .CommandPathPadding | cmdStyle }} {{ .Short }}
    {{- end }}
  {{- end }}
{{- end }}

{{- if .HasAvailableSubCommands }}
Use "{{ .CommandPath | cmdStyle }} [command] --help" for more information about a command.
{{- end }}
