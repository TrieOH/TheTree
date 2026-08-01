# Issue tracker: Forgejo (self-hosted)

Issues and PRDs for this repo live as Forgejo issues on the self-hosted instance at `https://git.trieoh.com` (org `TrieOH`, repo `TheTree`). The host is Gitea 1.22 (Forgejo-compatible API at `/api/v1`).

## Access

The API token lives at `~/.config/trieoh/forgejo.token` (chmod 600) — outside the repo, never commit it. Use it as `Authorization: token $(cat ~/.config/trieoh/forgejo.token)` on every request. There is no `gh`/`glab` CLI configured for this host; use `curl` against the Gitea REST API.

## Conventions

- **Create an issue**: `POST /api/v1/repos/TrieOH/TheTree/issues` with `{"title": ..., "body": ..., "labels": [<label-id>]}`.
- **Read an issue**: `GET /api/v1/repos/TrieOH/TheTree/issues/<number>` — `?comments=true` for comments.
- **List issues**: `GET /api/v1/repos/TrieOH/TheTree/issues?labels=<comma-separated-label-ids>`.
- **Comment on an issue**: `POST /api/v1/repos/TrieOH/TheTree/issues/<number>/comments` with `{"body": ...}`.
- **Apply / remove labels**: `PUT /api/v1/repos/TrieOH/TheTree/issues/<number>/labels` with `{"labels": [<label-ids>]}` (replaces the whole set); `DELETE /api/v1/repos/TrieOH/TheTree/issues/<number>/labels/<label-id>` to remove one.
- **Close**: `PATCH /api/v1/repos/TrieOH/TheTree/issues/<number>` with `{"state": "closed"}` — post the closing explanation as a comment first.
- **Labels**: `GET /api/v1/repos/TrieOH/TheTree/labels` to map names → ids.

The API is Gitea-shaped, not GitHub-shaped: paths start `/api/v1`, auth is a plain bearer-style `token` header, and issues and PRs are numbered in one sequence.

## Pull requests as a triage surface

**PRs as a request surface: no.** _(Set to `yes` if this repo treats external PRs as feature requests; `/triage` reads this flag.)_

## When a skill says "publish to the issue tracker"

Create a Forgejo issue via the API above.

## When a skill says "fetch the relevant ticket"

`GET /api/v1/repos/TrieOH/TheTree/issues/<number>` (with `?comments=true`).

## Wayfinding operations

Used by `/wayfinder`.

- **Map**: a single issue labelled `wayfinder:map`, holding the Notes / Decisions-so-far / Fog body.
- **Child ticket**: an issue carrying `Part of #<map>` at the top of its description and labels `wayfinder:<type>` (`research`/`prototype`/`grilling`/`task`).
- **Blocking**: text form only — a `Blocked by: #<n>, #<n>` line at the top of the description. A ticket is unblocked when every blocker is closed.
- **Frontier query**: list the map's children, drop any with an open `Blocked by` line or an assignee; first in map order wins.
- **Claim**: `PATCH` the issue with `{"assignees": ["<username>"]}`.
- **Resolve**: post a comment with the answer, close the issue, then append a context pointer to the map's Decisions-so-far.
