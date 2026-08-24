# Backend Wishlist

Status: **open — nothing built yet**. Date: 2026-08-24. Scope: `TrieOH/TheTree` — `api/{identityx,informd,payssage,univents}` · `lib/go` · `sdk/go` · CI workflows. Owner: **backend developer** (no infra component; the CI runner already exists).

From the backend dev point of view: language, dependencies, tooling, test hygiene. Go 1.27 shipped 2026-08-19 (five days ago) — two of these items are timed to it.

| # | Item | Effort | Why |
|---|------|--------|-----|
| 1 | **Upgrade Go to 1.27** (must) | M | Newest release; unlocks stdlib uuid (#2) + json/v2; repo is on 1.26.2–1.26.4 |
| 2 | **Replace google/uuid with stdlib `uuid`** (must) | M | RFC 9562 in stdlib since 1.27; drops a dep; ~600 files but mechanical |
| 3 | `-race` in CI integration tests | XS | `go test -count=1 ./...` runs today with no race detector |
| 4 | `go mod tidy` drift check in CI | XS | go.work.sum/go.sum drift goes uncaught between `just goup` runs |
| 5 | govulncheck in CI | S | Go-precise vuln scanning (call-graph aware) complements trivy's fs scan |
| 6 | Evaluate `encoding/json/v2` | L | 1.27's jsontext/json v2; bigger migration, defer until the harness needs it |
| 7 | Coverage reporting (optional) | S | No coverage signal anywhere today |
| 8 | OpenAPI spec linting (optional) | S | Specs are the contract (orval + oapi-codegen); lint them |
| 9 | Outbound HTTP timeouts (bug class) | S | Three `http.DefaultClient` calls with no timeout (MercadoPago, OAuth) |
| 10 | Rate limiting on auth endpoints | M | Zero rate limiting anywhere; login/token/OAuth are unthrottled |
| 11 | River + pgx pool metrics on `/metrics` | S | Job failures (webhook, email, key rotation) invisible except in logs |
| 12 | Span enrichment + error marking | M | Traces lack actor/project/provider context; business errors don't mark spans |
| 13 | Pagination across list endpoints | M | No limit/offset/cursor in any spec — list endpoints are unbounded |
| 14 | RFC 9457 Problem Details (off the fun envelope) | L | Envelope hides HTTP semantics; errors are non-standard |
| 15 | Hot-reload dev loop (Tilt / compose watch) | M | `just <svc>` rebuilds containers; no watch loop |
| 16 | Replace sdkkit in the Go SDKs (generate from spec) | M | Bespoke v0.0.2 kit used by 1 of 2 SDKs; couples public SDK to fun envelope |
| 17 | Resty everywhere an HTTP caller is needed | S | Already adopted in oauth_providers; the rest still hand-rolls clients |
| 18 | Dev seed data — `just seed` | S | Nothing seeds the dev DB; manual testing starts from empty postgres |
| 19 | Compose-scoped ops recipes — `just logs <svc>` / `just status` | XS | justfile only has host-wide `just ps` |
| 20 | `just nuke <svc>` — dev DB reset with confirmation | XS | Only blunt tool today is `down -v` (wipes everything) |
| 21 | Centralize tool versions | S | Pins duplicated across justfile / dagger / ci with "keep in sync" comments |
| 22 | Faster integration tests — shared testcontainer | M | `lib/go/testdb` pays container startup on every run |

---

## Must-haves

### 1. Upgrade Go to 1.27

**Why.** Go 1.27 (2026-08-19) is out and the repo is on 1.26.2 (payssage, lib/go) / 1.26.4 (identityx, informd, univents, go.work). 1.27 is what makes #2 possible (stdlib `uuid`) and ships `encoding/json/v2` (#6), `crypto/mldsa`, and the usual toolchain improvements. Sticking on 1.26.4 means the stdlib uuid migration can't start.

**Steps (one sweep, so nothing drifts):**
1. Bump `go.work` + all 5 `go.mod` files to `go 1.27`.
2. Bump the pinned toolchain in `.dagger/main.go` (`goVersion = "1.26.4"` → `1.27.x`), `ci.yml` (`setup-go` 1.26.4 → 1.27.x), and the `go-tools` Docker build image (rebuild with 1.27).
3. Bump pinned tool versions that predate 1.27: `sqlc` (1.31.1), `golangci-lint` (2.12.2), `gotestsum` (1.13.0), `oapi-codegen` (2.8.0) — check each for 1.27-compatible releases. Note: `sqlc` matters most — a newer sqlc may emit the stdlib uuid natively for pg `uuid` columns (see #2).
4. `just goup` to pull current dep releases (pgx, otel, zap, etc. all have 1.27-era versions).
5. `just generate` (sqlc + oapi + orval) and let the hermetic dagger CI validate the whole sweep.

**Gotcha.** After bumping the `go` directives, `GOTOOLCHAIN` auto-downloads 1.27 for anyone still on an older toolchain — that's the feature, not a bug: bump directives first, then `go work sync`.

### 2. Replace google/uuid with stdlib `uuid`

**Why.** Go 1.27 ships a native `uuid` package (RFC 9562, import path `uuid`, crypto-secure random). Dropping `github.com/google/uuid` removes a dep across all five modules. The DB type is unchanged (postgres `uuid` columns) — this is a pure Go type swap.

**Scale (verified):** 602 files, 2243 `uuid.UUID` type uses + 647 `uuid.New` + 28 `uuid.NewString` + 19 `uuid.Nil` + 16 `uuid.MustParse` + 6 `uuid.NewV7` + 5 `uuid.Parse` + 1 `uuid.Must(`.

**API differences that break compilation** (everything else is drop-in):

| google/uuid v1.6.0 | stdlib `uuid` | Impact |
|---|---|---|
| `uuid.New()` → `UUID` (v4) | `uuid.New()` → `UUID` (v4, identical) | ✓ drop-in |
| `uuid.UUID` `[16]byte`, comparable | same | ✓ drop-in |
| `uuid.MustParse`, `uuid.Parse`, `String()`, `MarshalText`/`UnmarshalText` | same | ✓ drop-in |
| `uuid.NewString()` → `string` | **none** — use `uuid.New().String()` | 28 sites, mechanical |
| `uuid.Nil` (var) | `uuid.Nil()` (func) | 19 sites, add `()` |
| `uuid.NewV7()` → `(UUID, error)` | `uuid.NewV7()` → `UUID` (no error) | 6 sites (all payssage), drop `err` |
| `uuid.Must(u, err)` | **none** | 1 site, drop wrapper |
| `Scan()`/`Value()` (sql.Scanner/driver.Valuer) | **none** | 0 direct sites — but see DB layer below |

**The DB layer — verified working, no adapter needed.** google/uuid worked with pgx because it implements `sql.Scanner`/`driver.Valuer`; stdlib `uuid.UUID` doesn't. But pgx v5.10 falls back to *underlying-type* plans for `[16]byte`-based types — I tested encode→scan round-trip with pgx v5.10 + stdlib `uuid.UUID` and it works (`*pgtype.underlyingTypeScanPlan` / `*pgtype.underlyingTypeEncodePlan`). The sqlc-generated models keep working unchanged after regeneration.

**sqlc.** All 4 `sqlc.yaml` explicitly override `db_type: "uuid"` → `github.com/google/uuid` (2 entries each). Change to `go_type: { import: "uuid", type: "UUID" }` (nullable + non-nullable variants), bump sqlc if a newer release knows the stdlib type natively, regenerate. Generated code is gitignored, so the swap is invisible in diffs beyond the override change.

**Careful spot.** The 6 `uuid.NewV7()` sites (payssage: idempotency keys for checkout/refund, intent IDs, webhook dispatch/receive) use V7's *sortable* property — stdlib `NewV7()` preserves it (48-bit ms timestamp + monotonicity guard, verified in source). Just drop the `err` from each call — don't switch them to `New()`/v4, that changes behavior (ordering).

**Steps:** sqlc overrides ×4 → regenerate → mechanical import swaps (codemod the 2243 type uses; hand-fix the ~50 sites in the table) → `go mod tidy` (uuid dep drops from all 5 modules) → full `just test` + `just lint` + CI.

---

## High-value

### 3. `-race` in CI integration tests

**Problem.** `ci.yml`'s `tests` job runs `go test -count=1 ./...` — no race detector. These are the DB-backed (testcontainers) tests, i.e. the ones with real concurrency, and they run on the only box where races bite. Local `just test` also omits `-race`.

**Fix.** `go test -race -count=1 ./...` in the `tests` job (and optionally in `just test`). Cost: slower runs (race builds) — acceptable for the integration-only job. Note: `-race` needs CGO; the runner image has gcc, fine.

**Where.** `.forgejo/workflows/ci.yml`, `justfile`.

**Effort.** XS.

### 4. `go mod tidy` drift check in CI

**Problem.** The workspace has 7 modules + `go.work.sum` (127KB). Drift between module `go.sum` files and `go.work.sum` accumulates between `just goup` runs and surfaces as confusing "updates to go.mod needed" errors at odd moments.

**Fix.** A CI step (or pre-commit): `go mod tidy && git diff --exit-code` — fails if tidy changes anything. One workflow in `ci.yml` (path-filtered to `**/go.mod`, `**/go.sum`, `go.work*`).

**Where.** `.forgejo/workflows/ci.yml` (+ `.husky/pre-commit` if wanted).

**Effort.** XS.

### 5. govulncheck in CI

**Problem.** trivy's fs scan (pre-push + PR workflow) is broad but not call-graph aware — it flags vulnerable packages even when the vulnerable code path is never reached.

**Fix.** `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` per module (or in the dagger `ci` step) as a complement: precise, only real reachable vulns. Add it to the `ci` job's hermetic pipeline.

**Where.** `.dagger/main.go` (`ci`), `.forgejo/workflows/ci.yml`.

**Effort.** S.

---

## Hardening

### 9. Outbound HTTP timeouts (bug class)

**Problem.** Three production outbound calls use `http.DefaultClient` — **no timeout**:
- `api/payssage/internal/features/providers/mercado_pago_provider/exchange_code.go:46` — MercadoPago token exchange (OAuth connect hangs forever if MP is slow)
- `lib/go/oauth/providers.go:118` — OAuth userinfo fetch (a hung Google/GitHub stalls login)
- (the healthcheck CLI's use is fine — deliberately `nolint`'d)

The good example exists: `deliver_webhook.go` builds its own client with a 10s timeout. The other two paths just don't.

**Fix.** A shared `lib/go/httpclient` (default 10s timeout + context) wired through the injected-client seams (CONTEXT.md: OAuth already takes an injected HTTP client — the seam exists, the default doesn't). Audit for any other outbound calls in the same sweep.

**Effort.** S. A live bug, not polish.

### 10. Rate limiting on auth endpoints

**Problem.** A global limiter **already exists** in the harness (`lib/go/httpserver/httpserver.go`): `mws.RateLimit(RPS: 400, Burst: 20, key: RemoteAddr)` — every route, every service. That stops floods, not credential stuffing: an attacker sustains 400 req/s per IP after the burst, and login/token-mint/action-token-redeem/OAuth connect get the same coarse limit as `/health`.

**Fix.** Auth-specific, tighter limits on the authn chain: per-IP + per-actor (e.g. 10/min on login, 60/min on refresh, action-token redeem tighter) — a second limiter ahead of the auth routes, keyed by `RemoteAddr` + actor id where known. `golang.org/x/time/rate` (already an indirect dep) or a second `mws.RateLimit` with a tighter config. Fail 429 with the RFC 9457 shape once #14 lands. (Note: this item was originally written as "zero rate limiting anywhere" — corrected 2026-08-24 after the harness's existing limiter was found; the gap is *specificity*, not presence.)

**Effort.** M.

### 11. River + pgx pool metrics on `/metrics`

**Problem.** River jobs run everywhere (webhook delivery, auth emails, key rotation, cleanups) and their failures are invisible except in logs; pool exhaustion is invisible too.

**Fix.** Export river queue depth / failed jobs / runs + pgx pool stats (`pgxpool.Stat`, and otelpgx is already a dep) on `/metrics` (the harness already serves prometheus). Pairs with infra wishlist #10 (latency alerts) for alerting on queue growth.

**Effort.** S.

### 12. Span enrichment + error marking

**Problem.** Traces carry route + status (otelhttp) but no actor/project/provider context, and business errors (4xx, swallowed errors) don't mark spans — trace search can't answer "which project saw this 500".

**Fix.** This is the **dev half of infra wishlist #11**: a `RecordError(ctx, err)` helper in `lib/go/telemetry` (`span.RecordError` + `SetStatus(codes.Error)`) called at caught-error sites; span attributes (`actor_id`, `project_id`, `provider`, job name) where available; wire the git tag into `semconv.ServiceVersion` (currently hardcoded `"dev"`).

**Effort.** M. Cross-ref: infra wishlist #11 (same work, observability framing).

### 17. Resty everywhere an HTTP caller is needed

**Problem.** resty v3 is **already adopted** in identityx's `oauth_providers` (an injected `resty.Client`, pinned at `v3.0.0-rc.3`) — but the rest of the codebase still hand-rolls clients: `lib/go/oauth/providers.go` and the MercadoPago `exchange_code.go` use naked `http.DefaultClient` (#9), `deliver_webhook.go` builds its own `http.Client`, and the IdentityX SDK sits on sdkkit's bespoke HTTP layer (#16).

**Fix.** A shared resty client factory in `lib/go/httpclient` (default timeout, retries with backoff, JSON marshaling, a structured-logging hook, otel transport so distributed traces keep flowing — the SDKs already do `otelhttp.NewTransport`) injected at every outbound seam: OAuth providers (done), MercadoPago provider, webhook delivery, SDK core. This **structurally kills the #9 naked-client class** — resty's defaults include timeouts — and adds retries for free. Note: `v3.0.0-rc.3` is fine for now; revisit when v3 stable lands.

**Effort.** S–M. Absorbs #9 (cross-ref).

---

## HTTP surface

### 13. Pagination across list endpoints

**Problem.** Zero pagination in any `api-spec.yml` (0 limit/offset/cursor refs across all four) — list endpoints return unbounded rows and will only get worse as data grows. The envelope even has a `PaginationMeta` field that's never populated.

**Fix.** Pick a scheme (**cursor** for high-churn lists like events/purchases, offset+limit is fine for the rest), standard params in the specs, a shared `lib/go` pagination helper (meta + `Link` header), and a harness-enforced default+max limit. Apply to every list operation.

**Effort.** M.

### 14. RFC 9457 Problem Details (move off the fun envelope)

**Problem.** Everything rides the `fun.Response` envelope (`{code, message, data, error, pagination, timestamp, module}`) — HTTP semantics are hidden inside a body, errors are non-standard, and every client couples to the envelope (the orval-mutator unwraps it, `ApiError` carries `.envelope`).

**Fix.** A deliberate architecture change: **plain JSON payloads on success + RFC 9457 `application/problem+json` for errors** (`type`/`title`/`status`/`detail`/`instance`) with proper HTTP status codes. Keep `fun` internally if useful, stop it at the HTTP boundary. Cross-cutting: specs → handlers → orval-mutator/`ApiError` → frontends, in one coordinated change with SDK version bumps. ADR-worthy.

**Effort.** L. Do it as one project with #13 (pagination rides the same HTTP-surface change).

### 16. Replace sdkkit in the Go SDKs

**Problem.** `github.com/MintzyG/sdkkit` (v0.0.2 — pre-1.0) is a bespoke base client used by **exactly one of the two Go SDKs**: IdentityX embeds `*sdkkit.Client`; Payssage is hand-written with a plain `http.Client` + otelhttp transport and works fine without it. The kit drags a private dependency and the fun envelope into a public SDK — and #14 (RFC 9457) removes the envelope it's coupled to. The internal TS clients are orval-generated from the spec; the public SDKs (Go + TS) are all hand-written and drift independently.

**Fix.** Two paths, pick by appetite:
- **(a) Generate the Go SDK clients from `api-spec.yml`** with oapi-codegen's client mode — the spec is already the contract for server bindings and orval; a generated Go client is always in sync, no hand-written plumbing, no kit. Drops sdkkit + fun + the hand-written client in one move. The deep fix — parity with the TS side's generated story.
- **(b) Swap sdkkit's HTTP layer for resty** (v3, already in the dep set) and keep the hand-written SDKs. Lighter; keeps the hand-written error model until #14 lands.

Either way, errors land on #14's problem+json shape (SDKError/APIError map to it).

**Effort.** M (a) / S (b). Cross-ref: #14 (envelope), #17 (resty).

---

## Dev loop

### 18. Dev seed data — `just seed`

**Problem.** Nothing seeds the dev DB; manual testing starts from an empty postgres and you hand-run the setup flows (org → project → actors → keys) before you can click through anything. The README quickstart stops at "the stack is up".

**Fix.** A `just seed` recipe that boots the stack, drives `/auth/setup` (the same bootstrap the staging plan uses), then creates a dev fixture: an org + project, actors with known credentials, an event/edition in univents, a TEST_MODE seller in payssage. Idempotent with a `--reset` variant. Turns "set up my sandbox" from ~15 minutes of API calls into one command.

**Effort.** S.

### 19. Compose-scoped ops recipes — `just logs <svc>` / `just status`

**Problem.** The justfile has only host-wide `just ps` (raw `docker ps`). Watching one service's logs is `docker compose logs -f <svc>` from memory; "is everything healthy" is a guess.

**Fix.** `just logs <svc>` (`docker compose logs -f <svc>`) and `just status` (`docker compose ps` with the health column). Two recipes, instant ergonomics.

**Effort.** XS.

### 20. `just nuke <svc>` — dev DB reset with confirmation

**Problem.** Resetting dev state means `docker compose down -v` (wipes **everything** — postgres + rustfs) or hand-dropping tables. There's no targeted "clean DB for this service" command, and the blunt tool is destructive.

**Fix.** `just nuke <svc>` → **confirms first** (`Are you sure? [y/N]`, default no — safety by default), then drops/recreates that service's schema (goose down-to-zero or `DROP SCHEMA ... CASCADE`) and re-applies migrations — leaves rustfs and the other services' data alone. Targeted, safe-by-default.

**Effort.** XS.

### 21. Centralize tool versions

**Problem.** goVersion / sqlc / golangci / gotestsum / oapi-codegen / trivy pins live in `justfile`, `.dagger/main.go`, and `ci.yml` — with hand-written "keep in sync" comments (e.g. `oapiCodegenVersion` vs the binary baked into the go-tools image). That's exactly the drift class that already bit (go-tools image, toolchain mismatches).

**Fix.** One source of truth — a `versions` file (or generated constants) consumed by justfile, dagger, and the CI workflow, plus a CI check that the copies agree. A tool bump becomes one edit.

**Effort.** S.

### 22. Faster integration tests — shared testcontainer

**Problem.** `lib/go/testdb` pays container startup on every run. *(Update 2026-08-24: the current design is one container per test binary, migrations once, truncate after each test — the per-test spawn was already killed. And CI has no test DB — a fresh runner spawns a container per run either way — so there is no "shared" DB to gain; the only real costs are the ~26 indirect modules testcontainers drags in.)*

**Fix.** If dependency hygiene ever starts to hurt: move the spawn to CI (`docker run` + `pg_isready` wait in a shared script), testdb connects via `TEST_DATABASE_URL`, dev gets a warm local DB via `just test-db`. Container lifecycle logic splits across three places — only worth it for module-count/supply-chain reasons, not speed.

**Effort.** M. Optional. **Superseded by simplify-wishlist #10** (which carries the corrected, CI-first design and the "keep testcontainers unless hygiene hurts" verdict).

---

## Stretch

### 6. Evaluate `encoding/json/v2`

**Problem/opportunity.** Go 1.27 graduates `encoding/json/v2` (+ `jsontext`): configurable options, stricter defaults, faster unmarshal (classic `encoding/json` is now backed by it). The harness's response envelope (`lib/go` fun) and all handlers serialize through `encoding/json`.

**Fix.** Don't migrate yet — the API surface differs (options-based, stricter edge-case behavior: duplicate keys, unknown fields, number handling). Worth a spike: run the test suite against `json/v2` to find behavioral deltas, then decide. This is a "when the harness needs it" item, not a now item — but it's the reason #1 is worth doing today.

**Effort.** L (spike first, then decide).

---

### 7. Coverage reporting

**Problem.** No coverage signal anywhere — `go test` runs but nobody sees a number.

**Fix.** `go test -coverprofile` in the `tests` job + a badge/summary in PRs. Low value for solo dev until the suite is meaningful; cheap to add later.

**Effort.** S.

### 8. OpenAPI spec linting

**Problem.** `api-spec.yml` files are the contract that orval (TS clients) and oapi-codegen (Go bindings) generate from — a spec typo ships as a broken client silently.

**Fix.** Spectral (or similar) on the 4 `api-spec.yml` files in CI: structural rules (required fields, consistent naming, no breaking `x-scope` typos). Cheap guard on the codegen pipeline.

**Effort.** S.

### 15. Hot-reload dev loop (Tilt or `docker compose watch`)

**Problem.** `just <svc>` rebuilds the container per change; the edit→run loop is slow.

**Fix.** **Tilt is a real option** — it orchestrates the compose stack (`docker_compose()` resources), file-watches → rebuild/restart, and gives a web dashboard of service status/logs. Cost: a `Tiltfile` + a running daemon, and it's k8s-oriented. **Lazier alternative that does the same core loop with zero new tooling:** `docker compose watch` (built into compose v2.22+, you're well past that) — watch a service's dir, rebuild/restart on change. Try compose watch first; reach for Tilt if the dashboard/multi-service control pulls.

**Effort.** S (compose watch) / M (Tilt).

---

## Explicitly out of scope

- **Pre-commit Go lint stage** — husky already lints/formats everyone repo-wide; a Go-specific stage adds nothing.

- **slog vs zap** — zap stays; OTLP log shipping (infra wishlist #9) is what matters, not the logger.
- **Dependency trimming beyond uuid** — the remaining direct deps (pgx, chi, river, otel, jwt, validator, jsonschema, aws-sdk for rustfs S3) are all load-bearing; no dead weight found.
- **Go version in the deploy image tags** — unrelated; infra's business.

## Suggested order

1. **#1 Go 1.27** + **#2 stdlib uuid** in one sweep (toolchain + pins + sqlc override + codemod + `just goup`).
2. **#9 outbound timeouts + #17 resty everywhere** — same window: the shared resty factory replaces the naked-client paths (#9's bug class) and the hand-built webhook client.
3. **#3 `-race`** + **#4 tidy check** + **#12 span enrichment** — small, ride the CI changes.
4. **#10 rate limiting** + **#11 river/pool metrics** — the security and ops gaps.
5. **#13 pagination + #14 RFC 9457 + #16 SDK generation** — one deliberate HTTP-surface project: ADR first, then specs → handlers → orval-mutator → generated Go SDKs (drops sdkkit + fun from the SDKs) → frontends (breaking; bump SDKs).
6. **#5 govulncheck** after the dust settles.
7. **#15 hot-reload** whenever — compose watch first, Tilt if the dashboard pulls.
8. **#6 json/v2 spike**; **#7/#8** optional whenever.
9. **Dev loop (as you go):** #18 seed, #19 logs/status, #20 nuke, #21 versions file, then #22 shared test DB.
