# Editor configuration

`ado` uses a text editor for editing tasks such as writing PR titles and
descriptions. By default, it uses the `$EDITOR` environment variable. If you
would like to use a different editor, add the `editor` value in `.ado.yml`.

Here is how to use Visual Studio Code:

```yaml
editor: "code --wait"
```

Config values for other common IDE/editors:

- JetBrains IDE: `idea --wait`, `rider --wait`, ...
- Visual Studio: `devenv /edit`

> **Note:** The editor command must be available on your system's `PATH`. If you
> get a "not recognized" error, either use the full path to the executable in
> the config, or add its directory to your `PATH` environment variable.