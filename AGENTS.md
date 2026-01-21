- Ask if unsure, unless explicitly asked to not do so.
- Always checking online official docs and related resources or up-to-date info.
- When making code changes:
  - Conform to exciting coding style and conventions
  - Use `mise build`, `mise test`, `mise lint` and other
    [mise tasks](https://mise.jdx.dev/tasks) instead of manually running
    commands.
  - Always aim to make minimal code change
- General coding styles
  - Prefer early-return-pattern over nested conditionals.
  - Try to keep the indentation level minimal.
  - Don't repeat yourself. Propose refactoring if needed.
  - Don't be complicated. Don't assume too much.
  - Ensure no lint issues. Don't use `//nolint`.
- Use modern UNIX CLI (e.g., rg, fd, ...) when possible;.
- Don't generate go test code unnecessarily.
- Generated files and temporary scripts should go into `.ai.dump` to prevent
  accidental commits.
- Use modern Golang features, e.g., generic, new std libs, ... when possible to
  keep code simple, concise.
