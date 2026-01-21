# `pipeline logs` views build logs

View logs from pipeline builds. Interactively select a pipeline, build, and job,
or specify them directly via arguments and flags.

## Examples

- `ado pipeline logs` - interactive selection
- `ado pl logs 12345` - by pipeline ID
- `ado pl logs buddy` - filter by keyword
- `ado pl logs ci deploy` - filter by multiple keywords (AND)
- `ado pl logs buddy -b 20240115.1` - specify build number
- `ado pl logs buddy -b 123 -s Deploy` - filter by stage name
- `ado pl logs buddy -n 100` - show last N lines
