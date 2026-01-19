# `pr create` submits a PR for the current branch

It uses ADO REST API to:

- Check if a PR from the current branch to the target branch already exists.
- If not, ask for confirmation before creating a new PR and ensure the PR is
  queued for policy evaluation (i.e., build, test, merge conflicts...).

It will ask for confirmation and push the branch to remote before submitting PR.
If the remote branch does not exist yet, it will also create that. It will open
`$EDITOR` to edit PR title and description.

The default target branch is the git default branch. Use `-t/--target` to
specify a different target branch.

All confirmation prompts and editing can be skipped by passing `-y/--yes` flag.

## PR title and description

`pr create` will use new commits between the current branch and target to
prepare PR title and description. It assumes that the commit messages follow the
convention below.

- The first line is the commit subject
- A blank line following the subject, to separate it with the body.
- The rest of the commit message is the body.

If the branch contains only one new commit compare the target branch, that
commit _subject_ and _body_ will be used as PR title and description.

If the branch contains more than one new commit:

- The PR title will be the branch name formatted with this template

  ```gotemplate
  {{.BranchName | replaceAll "/" "-"}}
  ```

- The PR description will be generated from commit messages using this template

  ```gotemplate
  {{range .Commits}}- {{.Subject}}{{end}}
  ```

Both these templates can be customized in the configuration file. For example:

```yaml
pull-request:
  create:
    pr_title_template: |
      {{ .BranchName | trimLeft "tai/"  |replaceAll "/" "-" }}
    pr_desc_template: |
      {{range .Commits}}
      - <details open="false"><summary>{{.Subject}}</summary>
          {{.Body}}
        </details>
      {{end}}
```

## Examples

- Submit PR to a specific branch: `ado pr create -t <branch_name>`
- Submit draft PR: `ado pr create --draft`
- Open browser after submitting: `ado pr create -b`
- Skip all prompt, accept the generated PR details: `ado pr create -y`
