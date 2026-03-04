# `pipeline builds` lists builds for a pipeline

List recent builds for a pipeline, showing build ID, status, result, branch,
and commit message.

## Examples

- `ado pl builds buddy` - list builds for pipeline matching "buddy"
- `ado pl builds -p 12345` - list builds by pipeline ID
- `ado pl builds buddy -n 20` - show 20 most recent builds
- `ado pl builds` - interactive pipeline selection
