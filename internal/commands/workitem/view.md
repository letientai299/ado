# View Work Item

Display detailed information about a single Azure DevOps work item.

Shows comprehensive information including:
- Title, ID, state, and type
- Assigned to, priority, area path, and iteration path
- Tags and parent work item
- Description (HTML stripped for readability)
- Relations (when using --relations flag)
- Created/modified dates and users

## Examples

```bash
# View by ID
ado wi view 12345

# View with relations (links to other items, commits, etc.)
ado wi view 12345 --relations
ado wi view 12345 -r

# Open in browser
ado wi view 12345 --browse
ado wi view 12345 -b

# Output formats
ado wi view 12345 -o yaml
ado wi view 12345 -o json

# Find by keyword (opens picker if multiple matches)
ado wi view "login bug"
```

## Flags

- `-b, --browse`: Open work item in browser
- `-r, --relations`: Include relations (links to other items, commits, PRs, etc.)
- `-o, --output`: Output format (simple, json, yaml)
- `-m, --mine`: When searching, only show your work items

## See Also

- [Work item documentation](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/about-work-items)
- [Work item fields](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/guidance/work-item-field)
