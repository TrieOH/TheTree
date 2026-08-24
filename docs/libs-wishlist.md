# Libs Wishlist — extraction into `lib/go`

Status: **open — nothing built yet**. Date: 2026-08-24. Scope: `TrieOH/TheTree` — `api/{identityx,informd,payssage,univents}` · `lib/go`. Owner: **backend developer**.

What's duplicated across the four services that should live in `lib/go`. Evidence-based: every item cites the actual files. The rule of thumb used here — extract when there are **≥2 copies of the same logic**, skip when the copies are just "same shape, different deps" (that's DI, not duplication).

| # | Item | Effort | Evidence |
|---|------|--------|----------|
| 1 | Shared service-authN middleware (`lib/go/authn`) | M | `auth.go` ×4 same shape (197 lines); `auth_dispatch.go` ×3 md5-identical; two conflicting verify paths |
| 2 | Webhook sign/verify (`lib/go/webhooks`) | XS–S | payssage `signPayload` HMAC + univents verify, hand-rolled on both sides |
| 3 | Config error reporting (`lib/go/config` — all env errors at once) | XS | env/v11 already aggregates; the one-line blob is the problem |
| 4 | Boot-skeleton residue audit (optional) | S | ~1267 lines of same-shape boot across 4 services — audit, don't rewrite |
| 5 | `lib/crypto` — absorb webhook HMAC | XS | payssage hand-rolls `crypto/hmac` while `crypto.HashHMACSHA256` exists |
| 6 | `lib/database` — activate `RegisterPoolMetrics` | XS | Exists, never called; backend #11's pool half |
| 7 | `lib/river` — metrics + periodic helper | S | Queue/failure metrics absent; periodic jobs hand-scheduled |
| 8 | `lib/email` — async river send-job | S–M | identityx wraps `Client` in its own river job; univents badges send via the same lib |
| 9 | `lib/oauth` — injectable HTTP client | S | `FetchGitHubEmail`/userinfo still on `http.DefaultClient` — the #9 bug lives in the lib |
| 10 | `lib/telemetry` — `RecordError` + attrs + `ServiceVersion` | S–M | Backend #12 / infra #11 enrichment point |
| 11 | `lib/authz` — chain-registration helper | XS | `app/auth_dispatch.go` ×3 do the same `NewResolver` wiring |

---

### 1. Shared service-authN middleware — `lib/go/authn`

**Problem.** Every downstream service (informd, payssage, univents) hand-assembles the same "verify identityx-issued access tokens + service API keys → dispatch to chain" middleware, in `internal/app/auth.go`:

- `idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)` + `fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)` — 4 files, 197 lines, same shape.
- The scope-dispatch chain (`app/auth_dispatch.go`) is **md5-identical across informd, payssage, univents** (28 lines each).
- **Two conflicting verify paths exist**: univents verifies via `lib/go/crypto` (JWKS); payssage and informd go through the identityx SDK client. Same job, two implementations.
- informd's `apiKeyHook` is a stub — `fun.ErrNotImplemented("api keys are not yet supported")`. The half-built feature is exactly what the shared seam would finish.

**Fix.** A `lib/go/authn` package with the seam exactly where the variability lives. The md5-identical part (scope dispatch chain) is library-side; the per-service part is a single **hook**, because each service already checks different things today and will check different things tomorrow:

```go
// identity := { Kind: jwt|api_key, Claims *idx.AccessClaims, Subject string }
authn.New(client, authn.Options{
    // per-service checks plug in here: univents' REQUIRED_VERIFIED_EMAIL
    // gate, platform-role checks, tenant gating, whatever comes next.
    Enrich: func(ctx context.Context, identity authn.Identity) (context.Context, error) {
        ...
    },
})
```

This preserves the existing `jwtHook` seam their `auth.go` files already pass to `fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)` — the extraction keeps that function, standardizes the verify/dispatch around it, unifies the two conflicting verify paths (univents via `lib/go/crypto`, payssage/informd via the SDK), and activates informd's stubbed API keys. Rate limiting (backend wishlist #10) hooks the same chain later.

**Effort.** M. Cross-ref: backend #10 (rate limiting), infra #11 (span enrichment — actor/project attributes land on the same chain).

### 2. Webhook sign/verify — `lib/go/webhooks`

**Problem.** payssage signs webhook deliveries with HMAC-SHA256 (`signPayload` in `services/webhooks/jobs/deliver.go`, header `X-Payssage-Signature`); univents verifies them with `PAYSSAGE_WEBHOOK_SECRET` (wired through `app/wire.go` into `services.NewOperations`). The sign and verify logic are hand-rolled on **both** sides of the same protocol.

**Fix.** A `lib/go/webhooks` package: `Sign(secret, payload)` + `Verify(secret, payload, sig)` (constant-time compare), header conventions shared. Both services use it; any future tenant receiving payssage webhooks gets it free. (MercadoPago's own signature scheme stays in the payssage provider — out of scope.)

**Effort.** XS–S. Small, self-contained, kills crypto logic from two services.

### 3. Config error reporting — `lib/go/config` (all env errors at once)

**Problem.** Reported as "if one required env is missing it short-circuits, so later errors aren't shown". Source-checked: **it doesn't short-circuit** — `caarlos0/env/v11@v11.4.1` (the pinned version) collects every field error into an `AggregateError` (`doParse` iterates all fields; `Unwrap() []error`, errors.Join-compatible). The real problem is rendering: `AggregateError.Error()` is one line (`env: A is required; B is required; C is required`) and `errx.Exit` dumps it as a single JSON field — technically all errors, practically unreadable, so it *feels* like only the first was reported.

**Fix.** A small `lib/go/config` helper that each service's `Load()` calls instead of `errx.Exit(err, ...)` directly:

```go
func Load() Config {
    var cfg Config
    if err := env.Parse(&cfg); err != nil {
        // prints one error per line (sorted, with the key), then exits
        config.ExitOnErrors(err)
    }
    return cfg
}
```

`ExitOnErrors` uses `errors.As(&env.AggregateError{})` to fan out every error (one per line, with its env key from the `env:` tag), so a config with five missing vars shows five lines. Also gives a typed hook (`errors.As` per error kind — `VarIsNotSetError` vs `ParseValueError`) for the day config validation wants to distinguish "missing" from "garbage".

**Effort.** XS. Cross-ref: backend wishlist #21 (tool versions) — the helper lives next to the env loader in `lib/go`.

### 4. Boot-skeleton residue audit (optional)

**Problem.** `internal/app/` still carries ~1267 lines of same-shape boot code across the four services (`app.go`, `router.go`, `run.go`, `constraints.go`, `wire.go`, `clients.go`). ADR-0001 already extracted the harness (`lib/httpserver`); what remains is DI assembly, route registration, config loading and constraint validation.

**Fix.** An **audit, not a rewrite**: after #1 lands (it removes the dispatch/auth residue), diff what's left service-by-service and move only the genuinely-identical slices into `lib/httpserver`. Deliberately **do not** abstract the irreducible parts — `clients.go` DI wiring (13–54 lines each, different deps) is per-service by nature, and over-abstracting the boot would undo ADR-0001's "thin services" win.

**Effort.** S (audit + small moves). Optional.

---

## Enrich existing libs

The other half of the story: the libs already exist, but the services bypass or under-use them. Every item below is a case of "the lib owns it, the service hand-rolls it anyway" or "the lib is one call short".

### 5. `lib/crypto` — absorb webhook HMAC

payssage's `deliver.go` imports `crypto/hmac` + `sha256` + `hex` directly to sign webhooks — while `lib/crypto.HashHMACSHA256` exists and isn't used there. Enrichment: payssage signs and univents verifies through `lib/crypto` (or the libs #2 webhooks helper built on it, with constant-time compare). The crypto is already in the lib; the services just don't call it.

**Effort.** XS.

### 6. `lib/database` — activate `RegisterPoolMetrics`

`RegisterPoolMetrics` exists in the lib and **no service calls it** (verified). Backend wishlist #11's pool half is already written — it needs one call per service at boot in `run.go`. The river half (queue depth/failures) is the only genuinely new piece.

**Effort.** XS. Cross-ref: backend #11.

### 7. `lib/river` — metrics + periodic helper

`lib/river` owns client/workers/register/migrate, but nothing exports river metrics (queue depth, failed jobs) onto `/metrics`, and periodic jobs (identityx's rotate_keys / cleanups) are hand-scheduled per service. Enrichment: a `Metrics()` hook wired into the harness's `/metrics`, and a periodic-registration helper so a schedule is one line.

**Effort.** S. Cross-ref: backend #11.

### 8. `lib/email` — async river send-job

identityx hand-rolls a river job around `lib/email.Client` (`jobs/send_auth_email.go`, `WorkerDefaults[emails.SendAuthEmailArgs]`); univents badges send through the same lib Client. Enrichment: `lib/email` owns a generic `SendJob` river worker + an `Enqueue(client, msg)` helper — neither service writes river plumbing for mail again.

**Effort.** S–M.

### 9. `lib/oauth` — injectable HTTP client

The #9 bug class lives *inside* the lib: `FetchGitHubEmail` and the userinfo fetch use `http.DefaultClient` (no timeout) at `providers.go:118`. Enrichment: the lib takes an injected client (defaulting to the shared resty factory from backend #17) — fixing the timeout issue for every consumer at once, and letting identityx's already-injected resty client be the one true path.

**Effort.** S. Cross-ref: backend #9/#17.

### 10. `lib/telemetry` — `RecordError` + attrs + `ServiceVersion`

`lib/telemetry` has Init/InitTracer/Log/StartSpan — nothing for error spans, span attributes, or real version labels. Enrichment: `RecordError(ctx, err)` (`span.RecordError` + `SetStatus(codes.Error)`), attribute helpers (actor/project/route), and `ServiceVersion` sourced from build info instead of the hardcoded `"dev"`. This is the dev half of infra #11 / backend #12, living where it belongs.

**Effort.** S–M. Cross-ref: infra #11, backend #12.

### 11. `lib/authz` — chain-registration helper

`app/auth_dispatch.go` ×3 (md5-identical) build `authz.NewResolver(spec.OpenAPISpec, authz.Primitives{...})` + resolve chains. The Resolver machinery is already in the lib — add a `ResolveChains(spec, primitives)` helper so that file shrinks to ~3 lines, folding into libs #1 (authn).

**Effort.** XS (part of #1).

---

## Explicitly not duplication (already good seams — don't touch)

- **Config loading** — all four services already use `caarlos0/env/v11` with the same tag pattern (`env:"KEY,required"`); only the structs differ, which is irreducible.
- **Email** — `lib/go/email` is the shared sender; `identityx/internal/emails` is service-specific templates + action tokens, correctly placed.
- **Object storage** — `lib/go/objectstorage` extracted; univents-only consumer today, but that's what a lib is for.
- **Repos/services facades** — `repos/repos.go`, `services/services.go` are thin aggregates by design (ADR-0002).
- **`xslices`** (81 uses), **`crypto`** (JWT verify), **`api_keys`** (issuer-only) — working as intended.
- **`globals`**, **`objectstorage`** — single-consumer libs (identityx setup gate; univents storage); fine as libs, no action.
- **`validator` / `jsonschema`** — used where needed (identityx schemas); no obvious bypass.
- **sqlc `models.go` / `db.go` ×4** — per-service schemas by nature; not duplication.

## Suggested order

1. **#1 authn** — the only item that changes behavior (fixes the two-path inconsistency, finishes informd API keys). **#11 folds in.**
2. **#2 webhooks + #5 crypto** — one small change: sign/verify through `lib/crypto`; do in the same window as backend #17 (resty).
3. **#3 config errors** — XS; do whenever, it's five minutes and improves every boot.
4. **#6 pool metrics + #7 river metrics** — together with backend #11; mostly activation + one new hook.
5. **#9 oauth client** — rides backend #17's resty factory.
6. **#10 telemetry enrichment** — the dev half of infra #11 / backend #12; do it with those.
7. **#8 email job** — whenever; it removes river plumbing from two services.
8. **#4 audit** — after #1, when the boot residue is smaller.
