# Create Work Item

Create a new Azure DevOps work item (Task, Bug, User Story, etc.).

## Examples

```bash
# Open editor to fill in all fields interactively
ado wi create

# Create a task with inline title
ado wi create --title "Fix login page performance"

# Create a bug
ado wi create -t Bug --title "Crash on empty input"

# Create a user story assigned to someone
ado wi create -t "User Story" --title "As a user, I want search" -A alice

# With description
ado wi create --title "Fix login" -d "The login page takes 5s to load"

# Specify area and iteration
ado wi create --title "New task" --area "MyProject\Backend" --iteration "MyProject\Sprint 5"

# Skip confirmation
ado wi create --title "Quick fix" -y

# Create and open in browser
ado wi create --title "New feature" -b
```

## Flags

- `--title`: Work item title (opens editor if omitted)
- `-d, --description`: Work item description
- `-t, --type`: Work item type (default: Task). Common types: Bug, Task, User Story, Feature, Epic
- `-A, --assignee`: Assign to a user by display name or email
- `--area`: Area path (e.g., Project\Team)
- `--iteration`: Iteration path (e.g., Project\Sprint 1)
- `-y, --yes`: Skip confirmation prompt
- `-b, --browse`: Open the created work item in browser

## Editor Mode

When `--title` is not provided, an editor opens with a structured template:

```
Title:
Type: Task
Assignee:
Area:
Iteration:

<!-- ado-wi-create: DO NOT REMOVE -->
Description:

```

Fill in the fields, save and close the editor to create the work item.
Any flags passed on the command line (e.g., `-t Bug`) pre-populate the template.
