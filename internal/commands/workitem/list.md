# List Work Items

Query and list Azure DevOps work items using WIQL (Work Item Query Language).

By default, lists work items assigned to you that are not in Closed, Done, or Removed states.

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

# Output formats
ado wi list -o yaml
ado wi list -o json

# Limit results
ado wi list -n 100
```

## Flags

- `-m, --mine`: Show only your work items (default behavior)
- `-a, --all`: Show all work items, not just yours
- `-t, --type`: Filter by work item type (Bug, Task, User Story, Feature, Epic)
- `-s, --state`: Filter by state (New, Active, Resolved, Closed, Done)
- `-n, --top`: Maximum number of work items to return (default: 50)
- `-o, --output`: Output format (simple, json, yaml)

## See Also

- [WIQL syntax](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)
- [Work item fields](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field)
