# ado – Azure DevOps CLI tool.

## Installation

**Linux/macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/letientai299/ado/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/letientai299/ado/main/scripts/install.ps1 | iex
```

Installs to `/usr/local/bin` (Unix) or `%USERPROFILE%\bin` (Windows) by default.
Use `-d <path>` to specify a custom directory. Use `--from-main` / `-FromMain`
to install the latest build from the main branch (requires [gh CLI][gh]).

[gh]: https://cli.github.com/
