# `pr analysis` — PR statistics for a date range (experimental)

> This command is experimental. Output format and available statistics may
> change in future releases without a deprecation notice.

Fetches all pull requests (active, completed, and abandoned) created within the
given date range and computes summary statistics.

## Flags

- `--from` — start of the date range (default: 30 days ago). Accepts RFC3339 or
  YYYY-MM-DD.
- `--to` — end of the date range (default: now). Accepts RFC3339 or YYYY-MM-DD.
- `--top` — how many contributors/reviewers to show in the "top N" tables
  (default: 5).
- `--output` / `-o` — output format: `simple` (default), `json`, `yaml`.

## Statistics

- Total PR count and breakdown by status (active / completed / abandoned)
- Draft ratio (drafts / total)
- Average and median active time for completed PRs (creation → close)
- Top N contributors ranked by PR count
- Top N reviewers ranked by number of PRs they participated in (vote ≠ 0)

## Examples

    ado pr analysis
    ado pr analysis --from 2025-01-01 --to 2025-03-01
    ado pr analysis --top 10 -o json
