# uptime-check

External uptime checks for the four public product APIs, as a Cloudflare
Worker — the outside observer. Everything else (Grafana, Beszel, the checks
themselves) dies with the box; this doesn't. Wishlist item #3.

**1 worker, 2 crons, 1 KV namespace, no containers on the VPS.**

## What it does

- **Every 5 min** (`*/5 * * * *`): probes each endpoint in `src/config.ts`.
  First attempt fails → up to 3 retries at 30s gaps, same invocation.
  - first try ok → **up**
  - recovered on retry → **degraded** — amber on the status page, silent,
    counts as up in averages
  - all attempts failed → **dead** — urgent ntfy ping (once on transition),
    recovery ping with downtime duration when it comes back
- **Daily** (`40 3 * * *`): compacts the retention ladder and rebuilds
  rollups:
  - 5-min raw points → kept **7 days**
  - hourly averages (degraded = up) → kept **30 days**
  - daily averages (avg of hourly) → kept **1 year**
- **Any request** to the worker: serves the status page from KV — current
  state per endpoint, a 7-day strip at 5-min resolution, 7d/30d/1y uptime %.
  Runs on Cloudflare's edge, so it answers even while the box is on fire.

Buckets are colored by the standard convention: 100% green, 0% red,
anything between amber. One KV write per 5-min run total (batched day key)
— ~300 writes/day, inside KV's free-tier 1k/day limit.

## Setup (one-time)

```sh
pnpm install

# 1. Create the KV namespace and paste the id into wrangler.jsonc
pnpm exec wrangler kv namespace create UPTIME

# 2. (optional) Serve the status page on a custom domain — uncomment
#    the routes block in wrangler.jsonc, or add it in the CF dashboard.

# 3. Secret: the ntfy topic for dead/recovery pings
pnpm exec wrangler secret put NTFY_URL   # e.g. https://ntfy.trieoh.com/uptime-alerts

# 4. Deploy
pnpm deploy
```

For local dev: `cp .dev.vars.example .dev.vars`, then
`pnpm dev --test-scheduled` and hit
`http://localhost:8787/__scheduled?cron=*%2F5+*+*+*+*` to trigger a check run.

## Deliberate limits

- Only the four product APIs are checked/published. Forgejo is dev-only,
  ntfy is the notifier itself — the status page is the fallback surface
  for both.
- A dead endpoint means ~1–2 min of retries before the phone buzzes
  (detection latency is retries × 30s + probe timeout). That's the
  noise/false-positive trade, not a bug.
- `fetch()` can't read TLS cert dates, so cert expiry stays with the #8
  server cron — don't fold it in here.
