# List Work Items

Query and list Azure DevOps work items using WIQL (Work Item Query Language).

By default, lists work items assigned to you that are not in Closed, Done, or
Removed states.

## Examples

```bash
# List your assigned work items (default)
ado wi list

# List with keywords filter
ado wi list "login" "bug"

# List all work items in the project
ado wi list --all

# Filter by type
ado wi list --type Bug
ado wi list -t Task
ado wi list -t "User Story"

# Filter by state
ado wi list --state Active
ado wi list -s New

# Combine filters
ado wi list --type Bug --state Active

# Filter by assignee
ado wi list --assignee alice
ado wi list -A "john.doe@company.com"

# Generate activity reports
ado wi list --report @Today-7    # Last 7 days
ado wi list --report 2026-02-01  # Since February 1st
ado wi list --report @Today-30 --assignee alice  # Alice's last 30 days

# Output formats
ado wi list -o yaml
ado wi list -o json

# Limit results
ado wi list -n 100
```

## Flags

- `-m, --mine`: Show only your work items (default behavior)
- `-a, --all`: Show all work items, not just yours
- `-A, --assignee`: Filter by assignee name or email (substring match); implies --all
- `-t, --type`: Filter by work item type (Bug, Task, User Story, Feature, Epic)
- `-s, --state`: Filter by state (New, Active, Resolved, Closed, Done)
- `--report`: Generate activity report for date range (e.g. 2026-01-01 or @Today-7)
- `-n, --top`: Maximum number of work items to return (default: 50)
- `-o, --output`: Output format (simple, json, yaml)

## Report Mode

The `--report` flag generates activity reports showing all work items (Closed,
Resolved) completed or worked on since the specified date. To format the report
output, define a custom output template (see below) named `report` and use
`-o report`.

**Date formats:**
- `@Today-N`: Relative days (e.g., `@Today-7` for last 7 days)
- `YYYY-MM-DD`: Absolute date (e.g., `2026-01-01`)

## Custom Output Templates

In the `ado.yml` config file, you can define custom output formats using [Go
templates][go_tpl]. Each template receives a list of work item views.

[go_tpl]: https://pkg.go.dev/text/template

**Available fields:** `ID`, `Title`, `State`, `Type`, `AssignedTo`,
`ChangedDate`, `ResolvedDate`, `ClosedDate`, `WebURL`.

**Extra template functions:** `groupByState` groups items by their state into a
`map[string][]WorkItemView`.

### Example: report template

```yaml
workitem:
  list:
    custom_output_templates:
      report: |
        {{- $byState := . | groupByState }}
        {{- range $state, $items := $byState }}
        ## {{ $state }} ({{ len $items }})
        {{- range $items }}
        - {{ .Title }}
        {{- end }}
        {{- end }}
```

Usage: `ado wi list --report @Today-7 -o report`

### Example: markdown links

```yaml
workitem:
  list:
    custom_output_templates:
      markdown: |
        {{ range . -}}
        - [#{{ .ID }} {{ .Title }}]({{ .WebURL }})
        {{ end }}
```

Usage: `ado wi list -o markdown`

You can also override the built-in `simple` format by defining a template with
that name.

## See Also

- [WIQL syntax](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)
- [Work item fields](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field)
