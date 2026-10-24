# Issue tracker: GitHub

Issues and PRDs for this repo live as GitHub issues on `TrieOH/TheTree` (https://github.com/TrieOH/TheTree). Use the `gh` CLI, already authenticated on dev machines.

## Conventions

- **Create an issue**: `gh issue create --title ... --body ... --label <name>`.
- **Read an issue**: `gh issue view <number>` — `--comments` for comments.
- **List issues**: `gh issue list --label <name>` (repeat `--label` to AND them).
- **Comment on an issue**: `gh issue comment <number> --body ...`.
- **Apply / remove labels**: `gh issue edit <number> --add-label <name>` / `--remove-label <name>`.
- **Close**: `gh issue close <number> --comment ...` — the comment carries the closing explanation.
- **Labels**: `gh label list`.

Issues and PRs share one number sequence.

## Pull requests as a triage surface

**PRs as a request surface: no.** _(Set to `yes` if this repo treats external PRs as feature requests; `/triage` reads this flag.)_

## When a skill says "publish to the issue tracker"

Create a GitHub issue with `gh issue create`.

## When a skill says "fetch the relevant ticket"

`gh issue view <number> --comments`.

## Wayfinding operations

Used by `/wayfinder`.

- **Map**: a single issue labelled `wayfinder:map`, holding the Notes / Decisions-so-far / Fog body.
- **Child ticket**: an issue carrying `Part of #<map>` at the top of its description and labels `wayfinder:<type>` (`research`/`prototype`/`grilling`/`task`).
- **Blocking**: text form only — a `Blocked by: #<n>, #<n>` line at the top of the description. A ticket is unblocked when every blocker is closed.
- **Frontier query**: list the map's children, drop any with an open `Blocked by` line or an assignee; first in map order wins.
- **Claim**: `gh issue edit <number> --add-assignee <username>`.
- **Resolve**: post a comment with the answer, close the issue, then append a context pointer to the map's Decisions-so-far.
