Call Azure DevOps REST APIs directly.

The api command provides a generic interface to Azure DevOps REST APIs,
automatically discovering available endpoints from the internal REST client.
This allows calling any API without writing custom command code.

## Endpoint Discovery

Endpoints are discovered automatically via reflection from the REST client.
The naming convention follows the client structure:

    Client.Group().Scope().Method() -> group.scope.method

For example:
- `Git().PRs(repo).List(ctx, query)` -> `git.prs.list`
- `Builds().ForProject(repo).ByID(ctx, id)` -> `builds.for_project.by_id`

## Parameters

Parameters are passed as key=value pairs after the endpoint path:

    ado api git.prs.by_id id=123

For struct parameters, use dot notation for nested fields:

    ado api git.prs.list list_query.top=10 list_query.search_criteria.status=active

Shorthand field names are supported when unambiguous:

    ado api git.prs.list top=10 status=active

Use `--describe` to see available parameters for an endpoint:

    ado api git.prs.list --describe

## Output

Results are output as JSON by default. Use `-o yaml` for YAML format.

## Tab Completion

Tab completion is fully supported:

- Endpoint paths: `ado api git.<TAB>` shows available endpoints
- Parameter names: `ado api git.prs.list --<TAB>` shows available flags
- Enum values: `ado api git.prs.list --status <TAB>` shows valid values
