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

The `--report` flag generates markdown-formatted activity reports suitable for copying to status updates or documentation. It shows all work items (Closed, Resolved) that were completed or worked on since the specified date.

**Date formats:**
- `@Today-N`: Relative days (e.g., `@Today-7` for last 7 days)
- `YYYY-MM-DD`: Absolute date (e.g., `2026-01-01`)

The output is grouped by state and includes a prompt suggestion for AI summarization.

## See Also

- [WIQL syntax](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)
- [Work item fields](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field)
