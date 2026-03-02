# ado – Azure DevOps CLI tool.

A fast, friendly CLI for working with Azure DevOps from your terminal. Manage
pull requests, pipelines, work items, and more without leaving the command line.

<!-- prettier-ignore-start -->
> [!NOTE]
> This is usable, but still mostly a WIP. CLI commands, flags, output format
> and behaviors might change without any prior notice.
<!-- prettier-ignore-end -->

<!-- toc -->

- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)

<!-- tocstop -->

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
Use `-d <path>` to specify a custom directory.

[gh]: https://cli.github.com/

## Usage

`ado` needs an access token to to fetch ADO REST API. You can provide a
[Personal Access Token][ado_pat] via the `ADO_PAT` env. However, it's
recommended to use [Entra ID Token][ado_entra] for enhanced security and
conveniency. Just make sure [azure CLI][az] is installed and authenticated.
`ado` will run `az account get-access-token` when needed.

[ado_pat]:
  https://learn.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate
[ado_entra]:
  https://learn.microsoft.com/en-us/azure/devops/integrate/get-started/authentication/entra
[az]: https://learn.microsoft.com/en-us/cli/azure/install-azure-cli

Assuming the token is acquirable, you can `cd` into one of your ADO repos and
try `ado pr list`.

If that command fails to acquires a token and listing your repo's PRs, check
`az account show` output and verify if the default account is in a different
tenant than your ADO account. If so, provide the right tenant via `--tenant`,
`ADO_TENANT`, or in `.ado.yml` config file (see `ado help config`).

Each command has detailed help. Run `ado help <command>` to learn more.

## Development

We use [mise][] to manage dev tools and tasks. See `mise tasks`.

[mise]: https://mise.jdx.dev

`ado` doesn't need to run inside an ADO repo dir. It just uses the git remote
URL to detect repo info for constructing API url. Use this minimal config
template to run `ado` against any of your ADO repos.:

```yaml
tenant: ...
repository:
  org: ...
  project: ...
  name: ...
```
