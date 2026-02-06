# Work Item Commands

Commands for listing and viewing Azure DevOps work items.

Work items are units of work tracked in Azure DevOps, such as bugs, tasks, user
stories, features, and epics.

## Examples

```bash
# List your assigned work items
ado workitem list

# List all active work items
ado workitem list --all

# View a specific work item
ado workitem view 12345

# View work item with YAML output
ado workitem view 12345 -o yaml
```

## See Also

- [Work Item documentation](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/about-work-items)
- [WIQL syntax](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)
