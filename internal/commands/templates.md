# Available Template Functions

ADO CLI uses [Go templates](https://golang.org/pkg/text/template/) for various output formatting.
For more information about Go template syntax (actions, pipelines, variables, etc.), see the Go
documentation. Below are the available template functions you can use in your config.

Utils (see `ado help pr list` for some example usages)

| Function     | Description                                              |
|:-------------|:---------------------------------------------------------|
| `markdown`   | Render markdown string                                   |
| `join`       | Join slice of strings with a separator: `join sep slice` |
| `indent`     | Indent all lines in the string: `indent num str`         |
| `trimSpace`  | Trim leading and trailing whitespace                     |
| `trimPrefix` | Trim prefix from string: `trimPrefix prefix str`         |
| `replaceAll` | Replace all occurrences: `replaceAll old new str`        |
| `tr`         | Replace using regex: `tr pattern replacementr str`       |

Coloring:

| Function    | Description                        |
|:------------|:-----------------------------------|
| `const`     | Apply style for constants          |
| `faint`     | Apply faint style                  |
| `warn`      | Apply warning style                |
| `error`     | Apply error style                  |
| `success`   | Apply success style                |
| `pending`   | Apply pending/waiting style        |
| `highlight` | Apply text matches highlight style |
| `h1`        | Apply H1 heading style             |
| `heading`   | Apply other headings style         |
| `person`    | Apply style for person names/email |
| `time`      | Apply style for time/dates         |
| `cmdStyle`  | Apply style for commands           |
