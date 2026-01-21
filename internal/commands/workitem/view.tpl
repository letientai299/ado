{{- $sep := "────────────────────────────────────────────────────────────────────────" -}}
{{h1 (printf "%s #%d" .Type .ID)}} {{.State | highlight}}
{{heading .Title}}

{{$sep}}
{{- if .AssignedTo}}
  {{const "Assigned To:"}} {{.AssignedTo | person}}
{{- end}}
{{- if .Priority}}
  {{const "Priority:"}} {{.Priority}}
{{- end}}
{{- if .AreaPath}}
  {{const "Area:"}} {{.AreaPath | faint}}
{{- end}}
{{- if .IterationPath}}
  {{const "Iteration:"}} {{.IterationPath | faint}}
{{- end}}
{{- if .Tags}}
  {{const "Tags:"}} {{join .Tags ", " | faint}}
{{- end}}
{{- if .ParentID}}
  {{const "Parent:"}} #{{.ParentID}}
{{- end}}
{{- if .Reason}}
  {{const "Reason:"}} {{.Reason}}
{{- end}}
{{$sep}}
{{- if .Description}}

{{const "Description:"}}
{{.Description | markdown}}
{{- end}}
{{- if .Relations}}

{{const "Relations:"}}
{{- range .Relations}}
  - {{.Type}}{{if .Name}}: {{.Name}}{{end}}
{{- end}}
{{- end}}

{{$sep}}
{{const "Created:"}} {{.CreatedDate | time}} by {{.CreatedBy | person}}
{{const "Changed:"}} {{.ChangedDate | time}} by {{.ChangedBy | person}}
{{- if .CommentCount}}
{{const "Comments:"}} {{.CommentCount}}
{{- end}}
{{faint .WebURL}}
