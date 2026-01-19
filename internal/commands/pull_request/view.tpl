{{.Title | h1 }}
{{.Description | trimSpace | markdown -}}
{{"PR Link:" | heading}} {{.WebURL}}
