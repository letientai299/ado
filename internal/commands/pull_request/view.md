# `pr view` show the details of a PR

If a numerical ID is provided, it tries to fetch that PR directly. Otherwise, it
searches for a PR whose title or description contains all the provided keywords.
The filtering logic is similar to `pr list`.

If multiple PRs match the keywords, an interactive picker will be shown to
select one.

## Examples

- Search and view a PR by title: `ado pr view "fix cve"`
- View PR by ID: `ado pr view 12345`
- Open PR in browser: `ado pr view 12345 -b`
