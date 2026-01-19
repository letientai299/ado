# Editor configuration

`ado` uses a text editor for editing tasks such as writing PR titles and
descriptions. By default, it uses the `$EDITOR` environment variable. If you
would like to use a different editor, add the `editor` value in `.ado.yml`.

The below example is for Visual Studio Code:

```yaml
editor: "code --wait"
```

Here are configuration values for other common IDE/editors:

- JetBrains IDE: `idea --wait`, `rider --wait`, ...
- Visual Studio: `devenv /edit`