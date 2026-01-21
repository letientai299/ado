{{- range . -}}
{{printf "%6d" .ID | highlight}} {{.Type | printf "%-12s" | faint}} {{.State | printf "%-10s"}} {{.Title}}
{{end -}}
