# `pr update` updates a PR for the current branch

This command updates a pull request's details.

By default, it will find all active PRs and ask you to select one. You can pass
a single numeric arg to be used as PR ID, or multiple args to be used as
filtering keywords. It supports the same filtering flags with `pr list`.

Use `-./--current-branch` to filter PRs whose source is the current branch. With
a typical development workflow where each branch is often associated with only
one PR, this flag helps skip the filtering process.

Use `-e/--edit` to open editor (see `ado help config editor`) for editing the
PR's title and description. The editor will show all commit messages as
reference below a marker line. Content at and after the marker line is ignored
when saving.

Use `-x/--execute` to execute one of the following actions on the PR.

- `approve`: approve the PR
- `reject`: reject the PR
- `resetVote`: reset your vote on the PR
- `complete`: mark PR as completed
- `publish`: publish the draft PR
- `draft`: mark PR as draft
- `abandon`: abandon the PR
- `reactivate`: reactivate an abandoned PR

Action availability depends on the PR status and user's permissions.

If `pr update` is run without any editing flags, it will iteractively ask for
each possible action.
