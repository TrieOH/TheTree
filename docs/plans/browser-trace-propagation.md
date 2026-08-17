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

## Config & dev toggle (first-class)

Tracing is **off by default** and enabled per environment. One config object
consumed by both browser init and the ingest server fn:

```ts
// front-core: tracing config (per-app, env-driven)
export interface TraceConfig {
  enabled: boolean          // master switch — when off, no spans, no ingest calls
  serviceName: string       // "univents-web" etc. (resource service.name)
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
| WS purchase socket | `ws purchase:<id>` | opt-in; needs backend traceparent query param |
| SSE inventory stream | `sse inventory:<id>` | opt-in; needs backend traceparent query param |

Start with the first five; add deliberately.

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
| `lib/ts/front-core/package.json` | 2 | add `sdk-trace-web`, `otlp-transformer` |
| `lib/ts/front-core/src/tracing/exporter.ts` (new) | 2 | custom exporter → ingest server fn |
| `lib/ts/front-core/src/tracing/ingest.ts` (new) | 2 | `ingestTracesServerFn` (validation, dev/prod URL) |
| `front/*/wrangler.jsonc` + secrets | 2 | `TRACES_INGEST_URL`, `TRACES_OTLP_USER/PASSWORD` |
| `../infra/caddy/Caddyfile` | — | **no change** |

No CORS changes, no Go changes, no zone.js.
