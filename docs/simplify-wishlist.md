# Simplify Wishlist — deletion over addition

Status: **open — nothing built yet**. Date: 2026-08-24. Scope: `TrieOH/TheTree` — `api/*` · `lib/go` · Dockerfiles. Owner: **backend developer**. Lens: ponytail — "the best code is the code never written". Every item names what to *delete* or *shrink*, with the honest tradeoff.

| # | Item | Effort | What it deletes |
|---|------|--------|-----------------|
| 1 | Delete dead `oauth.FetchGitHubEmail` | XS | 1 dead function (zero callers, verified) |
| 2 | Correct backend #10: rate limiting **exists** | XS | A wrong claim; re-scopes the real gap |
| 3 | Move single-consumer `lib/go/globals` into identityx | XS–S | 1 lib with exactly one consumer = indirection |
| 4 | Revisit ADR-0006 (consumerless encryption keys) | M | A whole module if approved — decision, not code |
| 5 | Decouple harness from `fun/middlewares` | M–L | Third-party middleware stack; 161 files import fun |
| 6 | Drop the global default tx runner (`SetDefaultRunner`) | M | Hidden package-global state (maybe not worth it — see note) |
| 7 | Boot glue dedup | XS–M | `strict.go` ×4 near-identical, Dockerfile healthcheck ×4, `auth_dispatch` ×3 (tracked in libs #1) |
| 8 | Audit 30 `nolint` + 36 TODO/FIXME | XS | Dead suppressions, stale comments |
| 9 | `aws-sdk-go-v2` → `minio-go` | M | 18 module lines for 2 S3 calls in one file |
| 10 | `testcontainers` → shared test DB | M | ~26 indirect modules; per-binary spawn already smart, this is deps + warm starts |
| 11 | Delete `lib/go/env` | XS | Zero importers repo-wide (verified) |
| 12 | Collapse `lib/go/utils` | XS | 2 funcs, 4 uses, all payssage |

---

### 1. Delete dead `oauth.FetchGitHubEmail`

**Verified:** `FetchGitHubEmail` (lib/go/oauth/providers.go:110) has **zero callers** repo-wide (lib, sdk, api, tests). GitHub email fetch isn't part of the userinfo flow — it's a leftover. Delete it and its `GitHubEmail` type if that becomes orphaned.

**Effort.** XS. One function gone.

### 2. Correct backend #10: rate limiting **exists**

**Finding:** the harness already wires a global limiter (`lib/go/httpserver/httpserver.go`): `mws.RateLimit(RPS: 400, Burst: 20, key: RemoteAddr)` — every route, every service, from `fun/middlewares`. Backend #10's "zero rate limiting anywhere" is wrong.

**What this changes:** the real gap is **auth-specific, tighter limits** (login/token/OAuth at 10/min per IP+actor), not a missing limiter. The 400rps/20-burst global stops floods, not credential stuffing. Re-scope backend #10 accordingly (fixed in the doc).

**Effort.** XS (doc correction) — and it makes #10's real work smaller.

### 3. Move single-consumer `lib/go/globals` into identityx

**Verified:** `globals` (MarkSetupComplete/SetupComplete — the "has setup run?" gate) is imported by **identityx only** (8 refs, all in that service). A lib with one consumer is indirection — move it into `api/identityx/internal/`.

**Counter (honest):** `objectstorage` is also single-consumer (univents) but is a real client abstraction with a future (informd storage) — **keep** it. `api_keys` is identityx-only today but exists to be the downstream verify path (libs #1) — **keep** it. `globals` is two tiny funcs with no expansion path — **move** it.

**Effort.** XS–S.

### 4. Revisit ADR-0006 — consumerless per-project encryption keys

**Verified:** the per-project RSA-4096 **encryption** key lifecycle (provision + rotate + sweep) has **no consumer** — CONTEXT confirms it: "OAuth secrets are sealed with the master key" (that's `crypto.EncryptPrivateKey`, a different path, and it IS used). ADR-0006 is a deliberate reservation for "future envelope-encryption".

**Ponytail view:** speculative complexity — a whole module provisioning and rotating keys nobody reads. Deleting it removes `internal/keys`' encryption half + its tests.

**The catch:** it's an ADR decision, not a code task. The item is **revisit the ADR**: if the reservation is still hypothetical after the RFC 9457 / SDK work (backend #14/#16), delete; if a consumer is coming, keep and drop this item. Flag for the next ADR review.

**Effort.** M (if approved: delete module + tests + rotation wiring).

### 5. Decouple the harness from `fun/middlewares`

**Finding:** the middleware stack (RealIP, RequestID, Logs, Metrics, Recover, Timeout, MaxBodySize, RateLimit, CORS) comes from **`fun/middlewares`** — third-party framework code — and 161 files across the repo import `fun`. The envelope removal (backend #14) is step 1 of decoupling; this is the rest.

**Ponytail view:** the harness (`lib/go/httpserver`) should own its stack natively — these are ~9 small middlewares, mostly stdlib-expressible (timeout, body-size limit, recover are ~20 lines each; CORS/rate-limit stay libs). Each one moved out of `fun` shrinks the framework surface.

**Honest effort:** M–L, incremental. Do NOT do this before #14; do it route-by-route. The payoff is a framework-free harness — the same reasoning that killed sdkkit (backend #16).

**Effort.** M–L. Cross-ref: backend #14.

### 6. Drop the global default tx runner (`SetDefaultRunner`)

**Finding:** `database.SetDefaultRunner(tx)` — a **package-level global** — is set once at boot in all four services; `RunTx`/`Queries` route through it. Hidden global state, but documented ("call once at startup") and working.

**Ponytail view:** globals are complexity — inject the runner instead. **But:** it's one call per service, zero bugs, and threading it through everything touches every repo call site. The lazy senior's actual answer: **leave it**, note the ceiling (`ponytail: package global, inject if tests ever need isolation`). Listed here so it's a conscious decision, not an accident.

**Effort.** M if done; **XS if consciously kept** (add the comment).

### 7. Boot glue dedup

- `app/strict.go` ×4 — 15 lines each, near-identical (differ only by service name): fold into the libs #4 boot audit.
- `auth_dispatch.go` ×3 md5-identical — already tracked in libs #1/#11.
- Dockerfiles: the `lib/cmd/healthcheck` build RUN line ×4 — one shared base stage kills the copy (XS).

**Effort.** XS–M, mostly already tracked elsewhere.

### 8. Audit 30 `nolint` + 36 TODO/FIXME

13 `gosec` suppressions (some are real workarounds, some stale), 8 `nilnil`, plus 36 TODO/FIXME across api/lib. A one-pass triage: delete dead suppressions, convert TODOs to issues or delete.

**Effort.** XS.

---

## Kill dependencies

### 9. `aws-sdk-go-v2` → `minio-go`

**Verified:** `lib/go/objectstorage` (the only consumer) uses exactly **2 S3 operations** + URL parsing: `PresignedPutObject` and `RemoveObject`. No GetObject, no listing, no buckets. Yet lib/go's go.mod carries **18 aws/smithy module lines** (config drags imds/sso/sts/oidc).

**Fix:** `minio-go` — one module, presigned-URL support. **Do NOT hand-roll SigV4**: presigned URLs are auth, and the ponytail rule is never simplify away security. minio-go keeps the signing, kills ~17 modules.

**Effort.** M.

### 10. `testcontainers` → `docker run` in CI + `TEST_DATABASE_URL` (optional — weakest kill)

**Verified:** testcontainers + its docker/moby/containerd/gopsutil chain ≈ **26 indirect modules** in lib/go alone. The current `lib/go/testdb` design is already good (one container per test *binary*, migrations once, truncate after each test).

**The honest framing:** the integration tests run **in CI on a fresh runner — there is no test DB**, and there never will be. A container is spawned fresh per run either way: testcontainers does it in-process, any alternative does it in the workflow. So this kill is **dependency/supply-chain hygiene only** — build-time + module-count surface, not runtime or speed. That makes it the weakest of the four kills, and optional.

**If done, the shape (CI-first, no "persistent test DB" premise):**
- `lib/go/testdb` keeps its API; backend: `TEST_DATABASE_URL` set → connect + reset (`DROP SCHEMA public CASCADE; CREATE SCHEMA public`) + migrate once per binary (mPath-keyed, as today); truncate-after-test unchanged. Unset → skip.
- A tiny `scripts/testdb.sh` (`docker run -d postgres:18-alpine` + `pg_isready` wait + print URL) used by both CI and dev.
- CI `tests` job (already installs Docker): run the script, export `TEST_DATABASE_URL`, `go test -p 1 -count=1 ./...` per service, `docker rm`. One extra step.
- Dev: `just test-db` — the same script keeps a long-lived local postgres for warm runs (the only place "shared" applies, optional).

**Cost to accept:** container lifecycle logic splits across workflow + script + testdb (three places vs testcontainers' one); per-package isolation → `-p 1`.

**Verdict:** keep testcontainers unless module/supply-chain hygiene starts to hurt — it's the standard, self-contained, and works in both environments with zero wiring.

**Effort.** M if done. Supersedes libs-wishlist #22.

### 11. Delete `lib/go/env`

**Verified:** zero importers repo-wide (api, lib, sdk, dagger). Config uses `caarlos0/env`; nothing uses `env.Get`. Delete the package.

**Effort.** XS.

### 12. Collapse `lib/go/utils`

**Verified:** 2 funcs (`isNil`, `map_to`), 4 uses — all in payssage (`mercado_pago_provider`, `intents`). Move into payssage (or inline).

**Effort.** XS.

---

## Deliberately NOT simplifying (the boundaries)

- **Action-token anti-replay** (identityx `tokens.ActionTokenManager` + persisted blacklist) — security property; never simplify away. Ponytail rule: security is not a simplification target.
- **ADR-0006 if kept** — same rule, deliberate reservation.
- **`univents/checkouts/checkout.go` (973 lines)** — splitting a big file ≠ simplifying; it's one coherent checkout flow.
- **`xslices` (84 uses), `validator`, `repos/services` facades** — legit, no stdlib equivalent / by-design (ADR-0002).
- **The `fun` envelope removal** — real, but it's backend #14, not this list.

## Suggested order

1. **#1 dead code + #2 rate-limit correction + #8 nolint audit** — a cleanup afternoon, zero risk.
2. **#11 env + #12 utils** — XS deletions, do right now.
3. **#3 globals move** — with the next identityx refactor.
4. **#4 ADR-0006 revisit** — schedule for the next ADR review, decide once.
5. **#9 aws → minio** — with the next objectstorage touch.
6. **#10 testcontainers → `docker run` in CI** — *optional*; only if module/supply-chain hygiene starts to hurt. Keep testcontainers otherwise.
7. **#7 boot glue** — rides the libs #4 audit.
8. **#5 fun decoupling** — behind backend #14; never ahead of it.
9. **#6 default runner** — consciously keep + comment, revisit if tests hurt.
