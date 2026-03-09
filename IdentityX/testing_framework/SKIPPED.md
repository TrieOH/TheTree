# Skipped Tests

This document tracks Go integration tests that were intentionally not ported to the TypeScript/Vitest framework, along with the reason for each skip.

---

## Key Rotation (`key-rotation`)

**Go test:** `testKeyRotation` → `t.Run("Key Rotation", ...)`

**Reason:** Both `GlobalKeyRotation` and `ProjectKeyRotation` require direct access to internal Go infrastructure that is not exposed via HTTP:

- `sqlc.New(suite.DB)` + `queries.RotateSigningKeysForGoAuth/Project()` — internal sqlc queries to trigger key rotation
- `crypto.GenerateEd25519()` / `crypto.Encrypt()` — internal Go crypto package
- `queries.CreateKeyPair()` — direct DB insert of a new signing key
- `suite.App.Keys.RevokeKey()` — internal service method, no HTTP equivalent
- `queries.DeleteExpiredRevokedKeys()` — internal DB cleanup query

The HTTP-observable behaviors (token still works after rotation, revoked key → 401) cannot be triggered without first invoking these internal operations.

**Resolution:** Skip entirely. Key rotation is an infrastructure/ops concern. If HTTP admin endpoints for rotation are added in the future, this test can be ported.

---

## Scopes — DB Constraint Tests (`scopes.test.ts`)

**Go tests:**
- `t.Run("CreateGlobalScopeError", ...)`
- `t.Run("CreateProjectRootScopeError", ...)`
- `t.Run("CreateInvalidScopeType", ...)`

**Reason:** These tests call `queries.CreateScope()` directly via sqlc with malformed inputs to assert that the database returns specific constraint violations (`unique_violation`, `check_violation`). They do not exercise any HTTP endpoint — they test DB-level invariants enforced by Postgres constraints.

The error checking (`errx.IsUniqueViolation`, `errx.IsCheckViolation`) is Go-specific and maps to raw `pgconn.PgError` codes. There is no HTTP surface to test these constraints against.

**Resolution:** Skip. These are DB schema correctness tests, not API contract tests. They should live as Go unit/integration tests alongside the schema migrations.