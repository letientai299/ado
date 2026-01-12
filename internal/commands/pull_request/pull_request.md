# `ado` CLI - Azure DevOps Command Line Tool

> "Efficiency is doing things right; effectiveness is doing the right things."
>
> - Peter Drucker

| Command    | Alias                  | Description                       |
| :--------- | :--------------------- | :-------------------------------- |
| `pr`       | `pull-request`, `pull` | Manage Azure DevOps pull requests |
| `pipeline` | `pp`                   | Manage Azure DevOps pipelines     |
| `doctor`   | -                      | Run prerequisite checks           |

```go
var ppRun = &cobra.Command{
  Use:     "run",
  Aliases: []string{"c"},
  Short:   "create a pull request",
}
```
