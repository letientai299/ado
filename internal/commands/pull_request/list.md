# `pr list` show PRs for the current repo

This provides a quick way to review or export info of all open PRs. You can
provide one or more keywords to filter pull requests by their title or
description. The search is case-insensitive.

## Custom output formats

In the `ado.yml` config file, you can set define custom output formats using [Go
templates][go_tpl]. The below example provides a `markdown` template.

[go_tpl]: https://pkg.go.dev/text/template

```yaml
pull-request:
  list:
    custom_output_templates:
      markdown: |
        {{ range $pr := . -}}
        - {{ if $pr.IsDraft }}DRAFT | {{end}}[{{$pr.Title}}]({{$pr.WebURL}})
        {{ end }}
```

Running `ado pr list -m -o markdown` will produce something like this.

```markdown
- [fix: CVE-2026-...](https://dev.azure.com/...)
- DRAFT | [chore: clean up, modernize ...](https://dev.azure.com/...)
```

## Examples

- List all active PRs: `ado pr list`
- List your PRs: `ado pr list -m`
- Search for PRs containing some keywords: `ado pr list fix`
- Use other output format PRs as JSON: `ado pr list -o json`