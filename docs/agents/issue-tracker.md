# Issue tracker: Linear

Issues and PRDs for this repo live in Linear, team **TrieOH** (key `TRI`, issues read `TRI-123`). Linear syncs them to GitHub issues on `TrieOH/TheTree`, so a GitHub issue can exist for the same work; Linear is the source of truth. Agents use the Linear MCP tools (`list_issues`, `get_issue`, `save_issue`, `list_comments`, `save_comment`, …).

## Conventions

- **Service**: each service is a Linear **project**: `IdentityX`, `Payssage`, `Univents`, `Informd`, plus `Infra`. Put an issue in the project it belongs to; the dev portal maps issues to services this way.
- **Create an issue**: `save_issue` with `team: "TRI"`, a `title`, a `description` and the `project`. Check `list_issues` with a `query` first, so you don't file a duplicate.
- **Read an issue**: `get_issue TRI-123`. It leaves out comments; read those with `list_comments`.
- **List issues**: `list_issues` with `team: "TRI"`, filtered by `project`, `label`, `state`, `assignee` or `cycle`.
- **Comment on an issue**: `save_comment` on the issue.
- **Apply / remove labels**: `save_issue` with `addLabels` / `removeLabels`. `labels` replaces the whole set, so don't use it to add one.
- **Close**: set the state to `Done` (or `Canceled`) with `save_issue`, after a comment that explains why.
- **Labels**: `list_issue_labels` with `team: "TRI"`. Type labels: `Bug`, `Feature`, `Improvement`.

## Pull requests as a triage surface

**PRs as a request surface: no.** _(Set to `yes` if this repo treats external PRs as feature requests; `/triage` reads this flag.)_

## When a skill says "publish to the issue tracker"

Create a Linear issue in team `TRI` with `save_issue`.

## When a skill says "fetch the relevant ticket"

`get_issue TRI-<n>`, then `list_comments` for the discussion.

## Wayfinding operations

Used by `/wayfinder`.

- **Map**: a single issue labelled `wayfinder:map`, holding the Notes / Decisions-so-far / Fog body.
- **Child ticket**: a sub-issue of the map (`parentId` set to the map) with label `wayfinder:<type>` (`research`/`prototype`/`grilling`/`task`).
- **Blocking**: Linear's `blocked by` relation between issues. A ticket is unblocked when every blocker is `Done` or `Canceled`.
- **Frontier query**: `list_issues` with `parentId` set to the map, drop any with an open blocker or an assignee; first in map order wins.
- **Claim**: `save_issue` with `assignee: "me"`.
- **Resolve**: post a comment with the answer, set the issue to `Done`, then append a context pointer to the map's Decisions-so-far.
