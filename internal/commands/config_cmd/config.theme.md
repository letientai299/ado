# Theme Configuration

`ado` supports customizable themes to personalize the terminal output
appearance. Themes control colors for various UI elements, including text,
errors, warnings, Markdown rendering, and syntax highlighting.

Without any configuration, `ado` pick between three built-in themes based on the
terminal's configuration: dark, light, or noTTY (text only, no color or special
codes, useful when piping output to other commands).

Themes can be configured in your `.ado.yml` (see `ado help config`). However,
the theme object is large, so it's recommended to load them from external files
using `include!` directive. See `etc/themes` for provided theme files.

```yaml
theme:
  include!: path/to/theme.yaml
```

To disable color (and other terminal rendering control codes), set `COLOR=never`
or `NO_COLOR=true`.
