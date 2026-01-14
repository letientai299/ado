# `ado` CLI - Azure DevOps Command Line Tool

> "Efficiency is doing things right; effectiveness is doing the right things."
>
> - Peter Drucker

## Commands Overview

| Command    | Alias                  | Description                       |
| :--------- | :--------------------- | :-------------------------------- |
| `pr`       | `pull-request`, `pull` | Manage Azure DevOps pull requests |
| `pipeline` | `pp`                   | Manage Azure DevOps pipelines     |
| `doctor`   | -                      | Run prerequisite checks           |

---

### Pull Request Commands

**Manage** pull _requests_ in your repository.

- `list` (alias: `ls`): List pull requests.
- `create` (alias: `c`): Create a new pull request.
- `update` (alias: `u`): Update an existing pull request.
- `browse` (alias: `view`, `v`): Open the pull request in your default web browser.

**Example: Creating a PR**

```bash
# hello world
ado pr create --title "Fix: bug in auth" --branch feature/auth-fix
```

### Pipeline Commands

Manage [and google](https://google.com) monitor your Azure DevOps pipelines.

- `list` (alias: `ls`): List available pipelines.
- `run` (alias: `c`): Trigger a pipeline run.
- `browse` (alias: `u`): View recent pipeline runs in the web interface.

**Example: Running a pipeline**

```bash
ado pipeline run --id 12345
```

1. `git` installation.
2. `az` (Azure CLI) installation.
3. Azure CLI authentication status.

---

For more information, visit the [official Azure DevOps Docs][ado-docs].

[ado-docs]: https://learn.microsoft.com/en-us/azure/devops/
