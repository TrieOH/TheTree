# ISSUE-0001: ToS continued-use consent is stamped at login only, not at token refresh

**Status:** open (known gap, accepted for now)
**Scope:** api/identityx — Terms of Service (`tos_acceptances` ledger)
**Introduced by:** feat/identityx-terms-of-service

## Context

The ToS feature records consent in `tos_acceptances` with a `source` of
`clickwrap` (registration box-check / accept endpoint) or `continued_use`
(first successful authentication after a version's `effective_at`).

`StampUse` — the continued-use stamp — is wired into exactly two touch
points:

- `authn.Login` (password login)
- OAuth callback for an *existing* identity (`updateExistingIdentity`)

Token **refreshes** (`POST /auth/refresh`) mint new sessions without
calling `StampUse`.

## The gap

A user who authenticates *before* v2's `effective_at` and then stays
continuously logged in across it (session kept alive by access-token
refresh) never triggers a login, so no `continued_use` row is recorded
for v2 — even though they are, contractually, using the service under
the new version (their frontend has been notified, continued use is
acceptance).

Consequence: the ledger can show a hole for such users ("last accepted:
v1") while they may already be bound by v2. The legal exposure is low —
the notice email is sent and persisted, and continued use is the
announced model — but the per-user evidence trail is incomplete until
their next fresh login.

## Why it was deferred

- Refresh is token plumbing, not an authentication event; stamping there
  dilutes the meaning of the ledger row ("how did they accept? → a
  background token renewal").
- The stamp-on-login path already covers the dominant case: access
  tokens are short-lived and any real re-authentication stamps.

## Options when this gets picked up

1. **Stamp in `authn.Refresh`** — same `StampUse` call, one-line risk;
   refreshes become consent events. Cheapest fix.
2. **Sweep job** — periodic River job stamps every actor of a project
   whose latest accepted version < current after `effective_at`. Catches
   truly dormant-but-logged-in users and closes the ledger for the whole
   population at once; touches more rows per cycle.
3. **Hard wall** — block refresh (or mint a "consent required" error) for
   actors on an outdated version; forces explicit re-acceptance.
   Strongest consent, frontend work required, contradicts the
   continued-use model.

Option 1 is the natural next step if the gap proves material.

## Workaround

None needed for correctness — logins self-heal the ledger. The gap is
detectable via `GET /projects/{id}/terms-of-service/acceptances`:
actors without a row for the current version past `effective_at` who are
actively refreshing are the gap population.
