# Delete Work Item

Delete an Azure DevOps work item by ID.

By default, the work item is sent to the **Recycle Bin** and can be restored
later from the Azure DevOps web UI.

Use `--destroy` to **permanently** delete a work item. This action is
irreversible.

## Examples

```bash
# Delete (moves to Recycle Bin)
ado wi delete 12345
ado wi rm 12345

# Skip confirmation
ado wi delete 12345 --yes
ado wi rm 12345 -y

# Permanently destroy (cannot be undone!)
ado wi delete 12345 --destroy
ado wi delete 12345 --destroy --yes
```

## Flags

- `--destroy`: Permanently delete the work item (cannot be recovered)
- `-y, --yes`: Skip confirmation prompt
