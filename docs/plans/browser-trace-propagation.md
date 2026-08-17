# Browser → API trace propagation — investigation & fix plan

Status: **investigated, not implemented**
Scope: browser (TanStack Start SPA on Cloudflare Workers) → Go APIs (univents, identityx, informd, payssage)

## Policy

**All traced traffic escapes through the BFF — safely and auth-protected.**
Browser trace context and browser span export flow exclusively through the
Cloudflare Worker (TanStack Start server functions). Requests that don't pass
through the BFF (direct/public fetches) are deliberately **untraced**.

**Traces are per-action, not per-page.** A SPA page is not a traceable unit —
a single page load spans hours of unrelated work, so "one trace per page view"
would produce giant, meaningless traces. Instead, **the front explicitly
chooses which specific actions/flows get traced** (checkout, publish event,
add-to-cart, login, …). Each traced action is its own trace; actions that call
the API hang the API trace (and its cross-service children) under them.

- No zone.js, no automatic instrumentation (no document-load / user-interaction
  / fetch auto-spans). Every span is placed deliberately at a chosen seam.
- A `session.id` attribute (optional) can group actions across a page view —
  but it is an attribute, **not a trace boundary**.

## Trace model

```
action:checkout                                ← new trace, new trace id
└── bff POST /checkouts                        ← child (traceparent threaded through BFF)
    └── univents: POST /checkouts/{id}         ← Go API root (extracts traceparent)
        ├── identityx: GET /.well-known/jwks.json
        └── payssage: POST /wallets/.../intents

action:add-to-cart                             ← separate trace, separate id
└── bff POST /cart/items
    └── univents: POST /cart/items
```

- The **outermost** traced action starts a new trace; **nested** traced calls
  (action → sub-action) are natural parent/child within it.
- Client-only actions (no API involved) export as standalone root traces.

## How traces reach the dashboard (data flow)

Two egress paths, one sink — both write into the **same Victoria Traces
instance**, joined by the shared trace ID:

```
browser span (withSpan)          Go API spans (otelhttp + Go SDK transport)
        │                                       │
        │ (1) BFF egress                        │ (2) server-side egress
        ▼                                       ▼
ingestTracesServerFn (Worker)        Go OTLP exporter (container, obs-net)
  dev:  127.0.0.1:10428                       victoria-traces:10428
  prod: traces.trieoh.com (basic auth)
        └───────────────┬───────────────────────┘
                        ▼
              Victoria Traces  ←  the sink
                        ▲
                        │ (read-only)
        ┌───────────────┴───────────────┐
        ▼                               ▼
  Grafana (grafana.trieoh.com)   Victoria Traces UI (traces.trieoh.com)
  provisioned datasource         basic-auth admin access
```

- The **BFF is the egress for browser-generated span data** (no other
  auth-protected path exists — `traces.trieoh.com` is basic-auth + CORS-less).
- **Go spans always exit server-side** via the container's OTLP exporter, for
  BFF-traced and public-traced calls alike.
- **Dashboards only read**: Grafana and the Victoria Traces UI query the store;
  they never see the ingestion path. The trace assembles because the API's root
  span carries the action span's trace ID.
- **Timing nuance**: API spans export immediately; the browser span lands when
  the batch flushes (time/size-based or `pagehide`) — for a few seconds the
  trace shows only the API part.

## Current request flows (the three edges)

1. **Public/direct**: browser → Go API directly (`publicFetcher`). **Untraced.**
2. **Authenticated via BFF (default)**: browser → Cloudflare Worker
   (`authenticatedProxyServerFn`) → Go API. **The traced path.**
3. **Auth ops**: browser → Worker server function → IdentityX. **Traced path.**

The front contains zero trace-context code today (verified by grep).

## Why OTel JS instead of hand-rolling

`@opentelemetry/api` already provides `propagation.TraceContext` (flags,
`tracestate`, validity) that hand-built `traceparent: 00-…` strings would
duplicate. ~4 KB, pure, no zone.js when used without automatic instrumentation.

Caveats: browser client instrumentation is officially "experimental / mostly
unspecified" (Browser SIG); and the traces backend is auth-gated (below), which
is exactly why export goes through the BFF.

## Key facts

- **No CORS changes needed.** `traceparent` never crosses origins from the
  browser: browser→Worker RPC is same-origin, Worker→API is server-side.
  (Also: `traceparent`/`tracestate` are not CORS-safelisted, so if anything
  ever goes cross-origin from the browser it would need the API
  `CORS_ALLOWED_HEADERS` updated — not the case here.)
- **Ingest targets** (`../infra`): prod = `https://traces.trieoh.com` (Caddy →
  `victoria-traces:10428`, behind basic auth — Worker supplies creds
  server-side); dev = `http://127.0.0.1:10428` (**published on localhost** in
  the infra compose) — reachable from the local vite dev server, no creds.
  This is what makes the dev toggle trivial.

## Tracing public/direct actions

Public actions (`publicFetcher` → Go API directly) are traceable the same way:
wrap the action with `withSpan`, and attach the resulting `traceparent` to the
direct fetch. Two facts make this cheap:

- **Preflight is amortized, not per-request.** `traceparent`/`tracestate` are
  not CORS-safelisted, so a direct cross-origin call carrying them preflights.
  The fun CORS middleware already defaults `Access-Control-Max-Age` to
  **10 minutes** (`lib/go/httpserver`, `mws.CORS`), so the preflight happens
  once per (origin, method, headers) combo and is cached.
- **Attachment is opt-in per action.** The fetcher only adds `traceparent` when
  the action was wrapped; untraced public calls stay simple requests (no
  preflight, no header).

Required change (the **one** place CORS config reappears in the whole plan):
add `traceparent,tracestate` to `CORS_ALLOWED_HEADERS` in
`api/{univents,identityx,informd,payssage}/.env` (and deployed envs). Without
this, a traced public call's preflight fails and the request is blocked — so
this is a hard dependency of tracing public actions. `traceparent` carries no
credentials — it's an opaque ID — so allowing it cross-origin is safe.

Note: SSE (`EventSource`) and WS cannot carry headers at all, so their public
streams still need the opt-in query-param touchpoint if traced (see
Non-goals).

## Config & dev toggle (first-class)

Tracing is **off by default** and enabled per environment. One config object
consumed by both browser init and the ingest server fn:

```ts
// front-core: tracing config (per-app, env-driven)
export interface TraceConfig {
  enabled: boolean          // master switch — when off, no spans, no ingest calls
  serviceName: string       // = backend app name, e.g. "univents" (see Service naming)
  ingestURL: string         // dev: http://127.0.0.1:10428/insert/opentelemetry/v1/traces
                            // prod: https://traces.trieoh.com/insert/opentelemetry/v1/traces
  ingestAuth?: { user: string; password: string }  // prod basic-auth creds, Worker secrets only
  sampleRatio?: number      // optional head-sampling for prod volume (default 1)
}
```

- **Browser** (`VITE_TRACING_ENABLED`): when disabled, the tracing module is a
  no-op (`withSpan` runs the callback with zero overhead, no spans created).
- **Server fn** (`TRACES_INGEST_URL` + `TRACES_OTLP_USER`/`TRACES_OTLP_PASSWORD`):
  the ingest target and creds never live in client code; dev points at
  localhost, prod at `traces.trieoh.com` via `wrangler secret`.
- The ingest server fn **no-ops when disabled or misconfigured** so dev can run
  without the observability stack.

## Service naming — front traces share the backend's services

Front span resources use the **same `service.name` as the backend app**
(`univents`, `identityx`, `informd`, `payssage` — matching the Go services'
`semconv.ServiceName(appName)`), **not** `-web` variants. Both write to the
same Victoria Traces store, so front traces land in the same Jaeger service
and the existing per-service "Recent Traces" panels with **zero dashboard
changes**.

Accepted trade-off: trace panels mix client and server spans under one
service name (a little crowding). Mitigate it with a distinguishing
attribute — the web SDK already sets `telemetry.sdk.language=webjs`
automatically; add `component: "web"` on front spans so you can filter to
client-only or server-only when needed.

Prometheus-based panels (requests/sec, latency, status codes, error rate)
are **unaffected** — the front sends no metrics; only the trace panels mix.

| Front app | `service.name` (same as backend) | Ingest target |
|---|---|---|
| univents | `univents` | victoria-traces (dev `127.0.0.1:10428` / prod `traces.trieoh.com`) |
| identityx | `identityx` | same |
| informd | `informd` | same |
| payssage | `payssage` | same |

## Phase 1 — action tracing + BFF propagation (no export yet)

### 1. Add `@opentelemetry/api` to `@trieoh/front-core`

### 2. Tracing module — `lib/ts/front-core/src/tracing/browser.ts` (new)

```ts
import { context, propagation, trace, type Attributes } from "@opentelemetry/api"

let cfg: TraceConfig = { enabled: false, serviceName: "front", ingestURL: "" }

export function initTracing(config: TraceConfig): void { cfg = config }

/** One traced action = one new trace (unless nested inside another action). */
export async function withSpan<T>(
  name: string,
  attrs: Attributes,
  fn: (traceparent: string) => Promise<T> | T,
): Promise<T> {
  if (!cfg.enabled) return fn("")
  const span = trace.getTracer(cfg.serviceName).startSpan(name, { attributes: attrs })
  const ctx = trace.setSpan(context.active(), span)
  const headers: Record<string, string> = {}
  propagation.inject(ctx, headers)   // captured synchronously — no async loss
  try {
    return await context.with(ctx, () => fn(headers.traceparent ?? ""))
  } finally {
    span.end()
  }
}
```

Nested `withSpan` calls inherit the active context → same trace, parent/child.
`fn` receives the `traceparent` to thread into the BFF proxy request (phase 1
already has the `headers` passthrough).

### 3. Thread traceparent through the BFF

- Browser: `withSpan("action:…", attrs, (traceparent) => proxyRequest({ path, headers: { traceparent } }))`
- BFF (`server.ts`): forward `traceparent` untouched on outbound `fetch` to the
  Go API, including auth paths (`observedFetch`).

### 4. Curated registry — the specific actions (starting set)

| Seam | Span name | Notes |
|---|---|---|
| Checkout flow | `action:checkout` | wraps create + confirm |
| Add to cart / update cart | `action:cart-update` | |
| Publish / edit event | `action:event-publish` | |
| Auth (login/register/refresh) | `action:auth` | via BFF |
| Image/variant upload | `action:upload` | presign (BFF) traced; link by object key |
| Public event/product listing | `action:view-event` | **direct** — needs CORS allowlist |
| WS purchase socket | `ws purchase:<id>` | opt-in; needs backend traceparent query param |
| SSE inventory stream | `sse inventory:<id>` | opt-in; needs backend traceparent query param |

Start with the first six; add deliberately.

## Phase 2 — export the curated spans through the BFF

1. **Provider only**: `@opentelemetry/sdk-trace-web` as `WebTracerProvider` +
   `BatchSpanProcessor` + custom exporter. **No automatic instrumentations.**
2. **Custom exporter → `ingestTracesServerFn`**: serializes spans
   (OTLP/JSON via `@opentelemetry/otlp-transformer`) and POSTs to the server
   fn, which forwards to `cfg.ingestURL` (dev: localhost:10428, prod:
   `traces.trieoh.com` + basic auth from Worker secrets).
3. **Ingest hardening** (it's a public Worker endpoint that can write traces):
   - payload size/shape validation (reject non-OTLP garbage),
   - same-origin check on the request,
   - rate limit / cap,
   - no-op when tracing disabled or ingest misconfigured.
4. **Flush policy**: time/size-based (e.g. every 5 s / 50 spans) +
   `visibilitychange` + `pagehide` — `sendBeacon` has a 64 KB cap, use
   `keepalive` fetch for larger batches.
5. **Caddy: no change.** Basic auth stays; only the Worker talks to
   `traces.trieoh.com`.

## Phase 3 — Worker hop instrumentation (optional)

`@cloudflare/workers-otel` / `otel-cf-workers` exporting server-side to
`cfg.ingestURL` (no CORS; `observability.enabled` already on). Seeds the
action traces server-side and makes the browser span real (removes the
"virtual parent" caveat from phases 1–2). Only if the BFF hop itself is worth
tracing.

## Non-goals / known limits

- Direct/public fetches untraced — by policy.
- No automatic instrumentation, no zone.js — the registry is the full surface.
- WS / SSE / OAuth-redirect tracing is opt-in and deferred (each needs a small
  backend touchpoint: `traceparent` query param on WS handshake / SSE endpoint,
  upload link via presign span, OAuth re-attach via `state`).
- Before phase 2, API traces show a dangling parent (the action span isn't
  exported yet) — expected, resolved by phase 2.
- No changes to the Go servers, external SDK packages, or metrics/logging.

## Verification

1. Unit: `withSpan` disabled = zero overhead no-op; enabled = new trace per
   action, nesting yields parent/child; `traceparent` parses to the action span.
2. BFF proxy test: `traceparent` forwarded on outbound fetch (mocked fetch).
3. Dev e2e: run against `127.0.0.1:10428` (no observability stack in prod), do
   a checkout, confirm one trace: `action:checkout` → univents →
   identityx/payssage.
4. Prod: same with `traces.trieoh.com`; confirm spans survive navigation.
5. Toggle test: tracing off → no ingest calls on the wire.
6. Public action test: wrapped `action:view-event` produces one trace
   (`action:view-event` → univents GET) across origins; untraced public calls
   show no `traceparent` and no preflight.

## Files touched

| File | Phase | Change |
|---|---|---|
| `lib/ts/front-core/package.json` | 1 | add `@opentelemetry/api` (catalog) |
| `lib/ts/front-core/src/tracing/browser.ts` (new) | 1 | config + `withSpan` action module |
| `lib/ts/front-core/src/tracing/config.ts` (new) | 1 | `TraceConfig` + env plumbing |
| `lib/ts/front-core/src/auth/tanstack/client.ts` | 1 | accept/forward `traceparent` on proxy requests |
| `lib/ts/front-core/src/auth/tanstack/server.ts` | 1 | forward `traceparent` on outbound fetches |
| `front/univents/src/env.ts` + app init | 1 | `VITE_TRACING_ENABLED`, wire `initTracing` |
| traced seams (`features/...`) | 1 | wrap chosen actions with `withSpan` |
| `api/{univents,identityx,informd,payssage}/.env` | 1 | `CORS_ALLOWED_HEADERS` += `traceparent,tracestate` (public actions only) |
| `lib/ts/front-core/package.json` | 2 | add `sdk-trace-web`, `otlp-transformer` |
| `lib/ts/front-core/src/tracing/exporter.ts` (new) | 2 | custom exporter → ingest server fn |
| `lib/ts/front-core/src/tracing/ingest.ts` (new) | 2 | `ingestTracesServerFn` (validation, dev/prod URL) |
| `front/*/wrangler.jsonc` + secrets | 2 | `TRACES_INGEST_URL`, `TRACES_OTLP_USER/PASSWORD` |
| `../infra/caddy/Caddyfile` | — | **no change** |

One config change: `CORS_ALLOWED_HEADERS` += `traceparent,tracestate` on the four APIs (only needed to trace public/direct actions). No Go code changes, no zone.js.
