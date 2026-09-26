# ISSUE-0002: ToS version can move between OAuth connect and callback

**Status:** open (known gap, accepted — do not fix unless it materializes)
**Scope:** api/identityx — OAuth registration × Terms of Service
**Introduced by:** feat/identityx-terms-of-service

## The gap

OAuth first-time registration gates consent at `connect` time: the
frontend shows the ToS and passes `accepted_tos=true`, and the callback
registers the actor recording **whatever version is current at callback
time**. If a project admin bumps the ToS between a user's connect and
their callback (state TTL bounds this to 10 minutes; the provider round
trip is usually seconds), the ledger row pins the new version while the
user actually read the old one.

The content hash is always internally consistent (hash of the version
recorded), so the ledger never contradicts itself — it just may describe
a document one revision newer than the one the user saw.

## Why it is accepted

- Probability: (ToS updates are rare, a few times a year) × (a bump
  landing inside a seconds-wide window). Near zero.
- Harm: bounded — the recorded version is the one in force, the user was
  in any case notified of the bump by the change email, and the notice
  window still applies to them.
- Fix cost: snapshot the accepted version into the one-time login state
  (`oauth_login_states` schema change + migration, or encode it into the
  opaque state token) and verify at callback. Touches the OAuth state
  model for a scenario that cannot realistically occur.

## If it ever needs fixing

Record `tos_version` on the login state at connect (only when the
project has ToS), and in `registerNewIdentity` accept only when the
current version equals the snapshotted one — a mismatch either re-asks
for consent (connect again) or records the snapshotted version with an
explicit note. Option A (re-ask) is the legally clean one.
