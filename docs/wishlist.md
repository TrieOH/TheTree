# Wishlist

One table for everything we want to do but haven't started. Backend, libs, CI and infra only. It replaces the old backend, frontend, infra, libs and simplify wishlists (they're in git history).

Rules:

- When an item gets a Linear issue, put the issue key in the last column. Linear holds the details from then on.
- When an item is done or dropped, delete its row. Don't strike it through.
- Size is a rough guess: XS < 1h, S ≈ half a day, M ≈ 1–3 days, L ≈ a week or more.

| # | Area | Item | Size | Why | Linear |
|---|------|------|------|-----|--------|
| 1 | Payssage | Fix the MercadoPago provider: refresh seller tokens, derive the idempotency key from the purchase | M | Expired tokens break checkout; retries can double-charge | TRI-13, TRI-14 |
| 2 | Backend | Timeouts on outbound HTTP: inject one client instead of `http.DefaultClient` | S | 3 calls can hang forever (MercadoPago `exchange_code.go`, `lib/go/oauth`, `httpserver/bootspan.go`) | TRI-15 |
| 3 | Backend | Email provider | M | Email is down everywhere; `lib/go/email` only speaks SMTP | TRI-12 |
| 4 | Univents | Move realtime to a Durable Object | M | The hub keeps sockets in container memory: open sockets keep the container awake, and it only works with `max_instances: 1` | TRI-16 |
| 5 | IdentityX | ToS consent edge cases (stamped at login not refresh; version can move during OAuth) | S | Consent records can be wrong | TRI-27, TRI-28 |
| 6 | Informd | Delete the duplicated namespaced routes | S | ~20 handlers marked "kill"; `authz.Form` already anchors via the form's namespace | TRI-17 |
| 7 | IdentityX, Informd | Delete the model constructors | S | ~8 constructors marked "kill"; build structs directly and validate at use sites | TRI-18 |
| 8 | Backend | Pagination on list endpoints | M | No `limit`/`cursor` in any of the 4 specs; SDK breaks are fine today | TRI-20 |
| 9 | Backend | Span enrichment, error marking, `ServiceVersion` | M | Traces lack actor/project/provider; business errors don't mark spans | TRI-23 |
| 10 | Backend | Cache the 80% (ADR-0020) | M | Fewer container wake-ups and DB hits | TRI-10 |
| 11 | Backend | Performance baselines per endpoint | M | Numbers to beat before optimizing cold starts and caching | TRI-22 |
| 12 | Infra | Latency alerts in Grafana Cloud | S | Traces already land there; nothing alerts on them | — |
| 13 | Backend | Circuit breakers on outbound calls | S | A slow MercadoPago or Payssage drags the caller down | TRI-21 |
| 14 | Libs | Webhook HMAC into `lib/go/crypto` | S | 3 hand-rolled `hmac.New` (payssage MP verify, payssage deliver, univents receive) | — |
| 15 | Libs | Delete dead `oauth.FetchGitHubEmail` | XS | Zero callers | — |
| 16 | Libs | `aws-sdk-go-v2` → `minio-go` | M | 18 module lines for the S3 calls in `lib/go/objectstorage` | — |
| 17 | Libs | Config: print env errors one per line | XS | `env` already collects them all, but they come out as one unreadable line | — |
| 18 | Backend | Revisit ADR-0006 (per-project encryption keys with no consumer) | M | Could delete a whole module; a decision first | — |
| 19 | Backend | Audit 30 `nolint` and 36 TODO/FIXME | XS | Dead suppressions and stale comments (many TODOs are rows 1, 6, 7) | — |
| 20 | CI | `-race`, `go mod tidy` drift check, govulncheck | S | Not in `.github/actions/go/test` today | TRI-25 |
| 21 | CI | One shared test database instead of testcontainers per run | M | Faster integration tests, fewer deps | TRI-24 |
| 22 | CI | OpenAPI spec linting | S | The specs are the contract for oapi-codegen and orval | — |
| 23 | CI | Coverage reporting | S | No coverage signal anywhere | — |
| 24 | Backend | Evaluate `encoding/json/v2` | L | Not urgent; do it when the harness needs it | — |
| 25 | Infra | Staging plan for Cloudflare | M | The old plan is VPS-shaped | TRI-26 |
| 26 | Infra | Decide on DB backups beyond Neon's restore window | S | Nothing lives off Neon today | — |
| 27 | Dev loop | `just seed` | S | Local testing starts from an empty DB | — |
| 28 | Dev loop | `just logs <svc>`, `just status`, `just nuke <svc>` | XS | Only host-wide `just ps` and `down -v` today | — |
| 29 | Dev loop | Hot reload (`docker compose watch`) | M | `just <svc>` rebuilds the container | — |
