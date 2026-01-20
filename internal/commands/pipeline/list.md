# `pipeline list` shows pipeline definitions

Lists all pipeline definitions for your repository. Filter by providing keywords
that match pipeline name or path (case-insensitive).

## Output formats

Use `-o` to change an output format: `simple` (default), `json`, `yaml`, or
custom templates defined in your `ado.yml` config (see `ado help pr list`).

## Examples

- `ado pipeline list` - list all pipelines
- `ado pl ls build` - filter by keyword
- `ado pipeline list -o json` - output as JSON
