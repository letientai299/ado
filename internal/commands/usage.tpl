{{- /*gotype:github.com/spf13/cobra.Command */ -}}
{{- if gt (len .Aliases) 0}}
{{ "Aliases:" | heading }} {{ range $i, $alias := .Aliases }}{{ if $i }}, {{ end }}{{ $alias | cmdStyle }}{{ end }}
{{- end }}

{{- if .HasExample }}
{{ "Examples:" | heading }}
{{ .Example }}
{{- end }}

{{- if .HasAvailableSubCommands }}
  {{- $cmds := .Commands }}
  {{- if eq (len .Groups) 0 }}

{{ "Available Commands:" | heading }}
    {{- range $cmds }}
      {{- if (or .IsAvailableCommand (eq .Name "help")) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
      {{- end }}
    {{- end }}
  {{- else }}
    {{- range $group := .Groups }}

{{ .Title | heading }}
      {{- range $cmds }}
        {{- if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help"))) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
        {{- end }}
      {{- end }}
    {{- end }}

    {{ if not .AllChildCommandsHaveGroup }}
{{ "Additional Commands:" | heading }}
      {{- range $cmds }}
        {{- if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help"))) }}
  {{ rpad .Name .NamePadding | cmdStyle }} {{ .Short }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end -}}

{{- renderFlags . -}}

{{ if .HasHelpSubCommands }}
{{ "Additional help topics:" | heading }}
  {{- range .Commands }}
    {{- if .IsAdditionalHelpTopicCommand }}
  {{ rpad .CommandPath .CommandPathPadding | cmdStyle }} {{ .Short }}
    {{- end }}
  {{- end }}
{{- end}}
{{ if .HasAvailableSubCommands }}
Use "{{ .CommandPath | cmdStyle }} [command] --help" for more information about a command.
{{- end -}}