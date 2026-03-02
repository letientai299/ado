# Work Items Report

{{- $byState := . | groupByState }}
{{- range $state, $items := $byState }}

## {{ $state }} ({{ len $items }})
{{- range $items }}
- {{ .Title }}
{{- end }}
{{- end }}

---
*Report generated with `ado workitem list --report`*

**For AI Summary:** Copy the work items above and ask: "Summarize this week's development progress into a brief status report suitable for stakeholders."