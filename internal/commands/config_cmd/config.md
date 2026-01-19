# `ado` Configuration

`ado` resolves its config from these sources and in their orders (the next
sources _might_ override previous ones):

1. Config file.
2. Environment variables. Only a few variables are supported.
3. Command-line flags
4. Auto-detect necessary info from git repo and Azure CLI.

## The config file

`ado` searches from the working dir up to git repo root, stop at the first file
matching the following pattern (in their order):

- `.ado.(yaml|yml)`
- `.config/ado.(yaml|yml)`. The `.config` dir is a emerging conventional
  location for storing per repo tooling config files, to not pollute the repo
  root.

The config file is in YAML format, with a special directive: `include!`. It
loads the content of another YAML file into the current node. Thus, it allows
externalizing and reusing complex config objects, such as `theme` or custom
`pr list` templates.

Use `etc/schemas/config.json` with your editor to have documentation in your IDE
when writing the config file.

## Environment Variables

ADO respects these environment variables and will use their values to override
the corresponding config values:

- `ADO_PAT`: Personal Access Token for authenticate to Azure DevOps API.
- `ADO_TENANT`: Azure tenant ID, used with `az` CLI to generate Microsoft Entra
  token, ignored if `ADO_PAT` is available. This could be auto-detected, but if
  you logged in to multiple tenants, and the default tenant is not the correct
  one to call ADO API, you can set this variable (or specify in the config
  file).
- `ADO_DEBUG`: Enable debug logging (set to "true" or "1").
- `EDITOR`: Default text editor command.

It's recommended to use `ADO_PAT` only if you can't use `az` CLI, because Entra
token is more secure.

## Auto-detect

Some commands need to work in the context of an ADO git repo, e.g.
`pull-request`, `pipeline`. For such commands, `ado` need to know the
organization, project and repo name. This info can be auto-detected by parsing
the git remote URL of the current repo.

Note that the auto-detection won't override existing values loaded from config
file or environment variables. This is to allow using `ado` outside of a target
ADO repo. For example, using `ado` to fetch and then analyze its pipeline log.