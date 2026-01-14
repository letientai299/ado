{{- /*gotype:github.com/spf13/cobra.Command */ -}}
{{ "Usage:" | headingStyle }}{{if .Runnable}}
  {{.UseLine | cmdStyle}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath | cmdStyle}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{ "Aliases:" | headingStyle }}
  {{.NameAndAliases | cmdStyle}}{{end}}{{if .HasExample}}

{{ "Examples:" | headingStyle }}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

{{ "Available Commands:" | headingStyle }}{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding | cmdStyle}} {{.Short }}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{ .Title | headingStyle }}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding | cmdStyle}} {{.Short }}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

{{ "Additional Commands:" | headingStyle }}{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding | cmdStyle}} {{.Short }}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{ "Flags:" | headingStyle }}{{range (flags .LocalFlags)}}
  {{flagName . | flagStyle}}
{{.Usage | wrap | indent 4}}{{end}}{{end}}{{if .HasAvailableInheritedFlags}}

{{ "Global Flags:" | headingStyle }}{{range (flags .InheritedFlags)}}
  {{flagName . | flagStyle}}
{{.Usage | wrap | indent 4}}{{end}}{{end}}{{if .HasHelpSubCommands}}

{{ "Additional help topics:" | headingStyle }}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding | cmdStyle}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath | cmdStyle}} [command] --help" for more information about a command.{{end -}}
