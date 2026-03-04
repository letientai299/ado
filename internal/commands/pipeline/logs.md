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
- `ado pl logs buddy --claude` - analyze logs with Claude AI
- `ado pl logs buddy -b 123 -j Deploy -n 200 --claude` - analyze specific job with Claude
- `ado pl logs buddy --claude -i` - start interactive Claude session

## Claude AI Analysis

Use `--claude` to pipe build logs to the Claude CLI for AI-powered analysis.
Claude will receive the log content along with build context (pipeline name,
build number, branch, job result) and print the analysis to the terminal.

Add `-i` (`--interactive`) to start an interactive chat session with Claude
instead, allowing you to ask follow-up questions about the logs.

Configure the Claude binary path:

- Config file: `claude: /path/to/claude`
- Environment: `ADO_CLAUDE=/path/to/claude`
- Default: `claude` (found on PATH)
