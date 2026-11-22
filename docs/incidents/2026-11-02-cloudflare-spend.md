---
title: Cloudflare Containers billed past the plan with no users
date: 2026-11-02
status: resolved
severity: medium
started: 2026-10-19
detected: 2026-11-02T11:27Z
mitigated: 2026-11-02T11:40Z
ended: 2026-11-02T17:50Z
updated: 2026-11-09
summary: Oversized containers kept awake by crons, cold-start calls and scanners cost about $23 over the $5 Workers Paid plan. Spend is now hard-capped inside the plan.
services: [identityx, payssage, univents, informd]
---

## What happened

A Cloudflare usage alert reported about $20 of metered usage on an account with no users, on top of the $5 Workers Paid plan.

The four backends (`identityx`, `payssage`, `univents`, `informd`) moved to Cloudflare Containers between Oct 19 and 23. Each ran up to 3 instances of the `basic` size (1 GiB of memory, 4 GB of disk) behind its Worker. Cloudflare bills a container's reserved memory and disk for every second it's awake, and bills the Durable Object that hosts it for the same seconds. The APIs use 40 to 55 MB. Crons, JWKS fetches on cold start, a chain of service calls at boot, and internet scanners kept the containers awake around the clock: 981 Durable Object hours in 9 days.

The first fix (one `lite` instance, a 2-minute sleep, fewer things waking them, and a spend guard) went live at 12:12 on Nov 2. A bug in the containers library then kept IdentityX awake for three more hours, and a second fix closed it at 17:50.

In short: we reserved 1 GiB × 3 instances × 4 apps for apps using about 40 MB, something always woke them before the 5-minute sleep, and the platform has no spending cap and reports overage a day late.

## Timeline

All times UTC.

| Time | What happened |
|---|---|
| Oct 19 22:20 | `identityx-api` container deployed: `basic`, `max_instances: 3`, `sleepAfter: 5m`, hourly and daily crons |
| Oct 23 20:37–21:56 | `payssage`, `informd` and `univents` deployed the same way. Payssage has a cron every 15 minutes |
| Oct 24 | Billing cycle starts. Durable Objects awake 92 hours that day |
| Oct 25 | Peak: 162 Durable Object hours in one day |
| Oct 26–31 | 77 to 147 hours a day. Once the 25 included GiB-hours are used, container memory costs $0.70 to $1.60 a day |
| Nov 1 | Durable Object duration passes the 400k GB-s included. The 34.2k over is billed as a whole 1M GB-s block: $12.50 in one day |
| Nov 2 11:27 | Usage alert email (about $20 metered). Investigation starts |
| Nov 2 11:40 | `max_instances` 3 → 1 on all four apps, through the API |
| Nov 2 12:08–12:12 | PR #68 deployed (`identityx/v0.40.0`, `payssage/v0.12.0`, `univents/v0.15.0`, `informd/v0.7.0`): `lite`, 1 instance, 2-minute sleep, fewer wakers, spend guard |
| Nov 2 15:00–17:50 | IdentityX stays awake with almost no traffic (the in-flight request leak, cause 4 below) |
| Nov 2 ~17:30 | Leak found in `@cloudflare/containers` 0.3.7: abandoned responses keep the in-flight count above zero |
| Nov 2 17:47–17:50 | PR #69 deployed (`identityx/v0.41.0`, `payssage/v0.12.1`, `univents/v0.15.1`, `informd/v0.7.1`). IdentityX restarted and asleep at 17:50:39 |
| Nov 2 ~17:58 | `univents-ui/v0.14.0` deployed, reading profiles through the query cache. Resolved |

## Numbers

Billable usage, Oct 24 to Nov 2, from the Cloudflare dashboard:

| Item | Used | Included | Cost |
|---|---|---|---|
| Durable Objects duration | 434.2k GB-s | 400k GB-s | $12.50 |
| Container memory | 3.38M GiB-s | 25 GiB-h (90k) | $8.23 |
| Container disk | 13.53M GB-s | 200 GB-h | $0.90 |
| Container vCPU | 38.02k vCPU-s | 375 vCPU-min | $0.31 |
| Everything else (Workers, R2, D1, KV, Queues, DO requests) | | | $0.00 |
| **Total metered** | | | **$21.93** |

- **Cycle total:** about $28, which is $21.93 metered, the $5 plan, and a few cents before the fix. The dashboard projected $68, but it treated the one-off $12.50 block as a daily cost.
- **Who was awake** (Durable Object hours, Oct 24 to Nov 2): IdentityX 399, Univents 274, Payssage 191, Informd 118. That's 981 hours, about four containers awake nonstop.
- **What the apps use:** 36 to 48 MB of memory on average, about 60 MB of disk, CPU near idle. IdentityX peaked once at 295 MB.
- **Workers was never the cost:** about 190k requests and 1.1M CPU-ms, against 10M and 30M included.
- **R2 was a misread:** 69 MB stored and about 1.7k operations, all free. The $12 line was Durable Objects.
- **The Durable Object math:** 981 h × 3600 s × 0.125 GB (every awake one is billed as 128 MB) ≈ 441k GB-s, close to the dashboard's 434k.

## Why

1. **How containers are billed vs. how big our apps are.** Containers bill reserved memory and disk for every second awake, plus Durable Object time for the same seconds. `basic` reserved 1 GiB for apps using about 40 MB, 25 times too much. `getRandom(env.X, 3)` spread requests over three instances, so a trickle of traffic woke all three, and each stayed up for its 5-minute `sleepAfter`.
2. **Something always woke them.**
   - Payssage's cron ran every 15 minutes (96 wakes a day) and IdentityX's cleanups every hour (24 a day). Each wake meant at least 5 minutes awake.
   - Every other API fetched IdentityX's `/.well-known/jwks.json` on cold start: about 430 requests in 3 days.
   - Univents' boot called `GetWallet` on Payssage, and Payssage's API key check called IdentityX's `/auth/introspect`. One cold start became three (55 calls on each side).
   - Scanners probing `/`, `/.env`, `*.php` and `gcp-key.json` reached the containers through the Workers.
3. **No hard cap.** Cloudflare has no spending limit for usage products. Budget alerts are email only and arrive the day after the threshold. Durable Object overage is billed in whole 1M GB-s blocks, so the first GB-s over cost $12.50.
4. **The in-flight request leak,** found after the first fix. `@cloudflare/containers` 0.3.7 counts a request as in flight until its response body is fully read, and a container never sleeps while that count is above zero. On the slower `lite` cold start, the edge gave up on JWKS requests (10 × 504). Nobody read those responses, so IdentityX stayed awake from 15:00 to 17:50 with no real traffic. The spend guard also counted too much: the library reports some stops late, and the guard counted the sleeping time as awake (Univents: 889 GiB-s counted vs. about 169 real).

## Evidence

- Cloudflare billable usage (Oct 24 to Nov 2), and the Nov 1 daily breakdown with the $12.50 block
- GraphQL Analytics: `durableObjectsPeriodicGroups` (awake hours per app and day), `containersUsageAdaptiveGroups` (reserved memory per 15 minutes), `containersMetricsAdaptiveGroups` (real memory and CPU), `workersInvocationsAdaptive`, `httpRequestsAdaptiveGroups` (which paths woke them), and the `r2*` datasets
- Containers API: app config (`max_instances`, instance size) and `instances-v2` (IdentityX running since 16:43:54 with no request after 17:16)
- Library source: `@cloudflare/containers` 0.3.7 `container.js`. `isActivityExpired()` renews while `inflightRequests > 0`, and the count only drops after `res.body.pipeTo()` settles
- Code: `lib/go/backend/identityx` (introspect), `api/univents/internal/app/clients.go` (boot wallet check), `sdk/go/IdentityX/tokens.go` (JWKS fetch)
- KV `container-budget`: awake GiB-seconds per app since the guard went live

## What we changed

| Status | Change |
|---|---|
| Done | **Cap instances right away.** `max_instances` 3 → 1 on all four apps at 11:40, through the Containers API, before any code change. |
| Done | **Smaller containers that sleep sooner** (PR #68). `api/*/cloudflare`: `instance_type` `basic` → `lite` (1/16 vCPU, 256 MiB, 2 GB), `INSTANCES = 1`, `sleepAfter` 5m → 2m, and `GOMEMLIMIT=200MiB` so Go collects before the 256 MiB limit. CPU-heavy work like password hashing is slower on 1/16 vCPU. |
| Done | **Fewer things waking them** (PRs #68 and #69). IdentityX crons once a day at 03:00 (cleanups and key rotation) instead of hourly. Payssage redrive every 6 hours instead of 15 minutes. Univents' boot wallet check off on Cloudflare (`PAYSSAGE_VERIFY_WALLET_ON_BOOT=false`). Workers answer scanner paths with 404 without waking a container (`isNoise()`). JWKS served from the edge cache and from KV for an hour, dropped right after `keys.rotate`. Access tokens last 1 hour instead of 15 minutes, so fewer refreshes, but a revoked token stays valid for up to an hour. The Univents front reads attendee profiles through the query cache. |
| Done | **A hard spend cap** in `lib/ts/container-budget`. `BudgetedContainer` adds up each container's awake GiB-seconds in KV `container-budget`. At 90% of the plan's 25 GiB-hours in a billing cycle (starting on the 24th), the APIs answer 503 instead of starting a container, until the next cycle. Spend can't leave the $5 plan, but the APIs go down for the rest of the cycle if the budget runs out. |
| Done | **Idle watchdog and correct counting** (PR #69). A container with no new request for 5 minutes is stopped, whatever the library's in-flight count says. Late stop reports are clamped to the last request plus that window, so sleeping time is never counted. Univents' WebSocket and SSE connections drop after 5 idle minutes and reconnect on next use. |
| Done | **Fixed the pre-commit hook.** Line 56 of `.husky/pre-commit` redirected `engine_targets` into the file it was writing, so every commit failed with "input file is output file". |
| Watch | **Check the guard against Cloudflare.** After a few days, compare KV `container-budget` with `containersUsageAdaptiveGroups`. Idle days should be near zero, and active use about 1 to 1.3 container-hours per hour. The plan covers about 90 `lite` hours a month, about 4.5 hours per weekday. |
| Later | **Faster cold starts on `lite`.** IdentityX's boot runs `keys.EnsureAll`, with one Neon round trip after another per project, before taking traffic. Slow boots caused the 504s behind cause 4, so less boot work means lower latency and less awake time. |
| Later | **Show the budget.** Expose spent and limit from KV (a small endpoint or a Grafana panel) and warn at 50% and 75%, so a cut-off on the 20th isn't a surprise. Lower the Cloudflare budget alert to $1 as a second signal. |
| Later | **Revisit the platform if usage grows.** Containers bill reserved memory, and `lite` (256 MiB) is the smallest size. Past about 90 hours a month, a fixed-price host would run all four APIs at about 40 MB each for less. |

## Follow-up: day 1 (2026-11-03)

The first 23 hours after the fix, Nov 2 18:00 to Nov 3 17:00 UTC. The fix holds. The four containers were awake 5.8 Durable Object hours in total, against about 98 hours a day before (−94%). Container memory dropped from about 338k to 3.1k GiB-s a day (−99%), under a cent a day at overage rates. No new block is coming: the one bought on Nov 1 covers up to 1.4M GB-s, and the day used about 2.7k.

| App | DO awake | Container awake | Awake in | Why |
|---|---|---|---|---|
| identityx | 230 min | 127 min | 13 of 23 hours | Other APIs' cold starts (JWKS), the 03:00 cron, real use |
| univents | 73 min | 59 min | 6 hours | Real use (08:00, 12:00–13:00) |
| payssage | 30 min | 16 min | 6 hours | The 6-hourly redrive cron, use |
| informd | 15 min | 5 min | 1 hour | One wake |

- **Workers:** about 4.2k requests, 0 errors. All four apps on `lite`, 1 instance, PR #69 images.
- **Spend guard:** 5.56k GiB-s counted since it went live, about 1.5 of the 22.5 GiB-h cut-off. 4.1k of it is IdentityX, including the 15:00–17:50 leak.
- **Headroom was tight.** At the day's 0.86 GiB-h, the cycle would end (Nov 24) near 19.6 of 22.5 GiB-h. A few busy days would have tripped the guard and taken the APIs down until the 24th.

What still woke IdentityX:

- **The JWKS copy in KV never stuck.** A day later KV had no `jwks:` key. The 1-hour TTL ran out between wakes, and on a `lite` cold start the caller gave up (499/504) before the fill, so nothing was written.
- **So every cold start still chained into IdentityX.** Payssage's redrive at 00:00 and 06:00 woke Payssage, and its JWKS fetch woke IdentityX. Each of those wakes lasted 4 to 6 minutes, because the abandoned JWKS request left the container to the 5-minute watchdog instead of the 2-minute sleep.
- **Scanners are mostly handled.** Hundreds of `.env` and `.php` probes got a 404 from the Workers. A few paths without extensions (`/info`, `/api/status`) still reach a container, rarely enough to leave for now.

| Status | Change |
|---|---|
| PR #70 | **Keep the JWKS until keys can change.** The KV copy lives 7 days and is dropped only when keys can change: container start (boot `EnsureAll`), `keys.rotate` and project creation. The fill runs in `waitUntil`, so it lands even if the caller times out. Other APIs' cold starts stop waking IdentityX. |
| PR #70 | **Skip the Payssage redrive on idle days.** The Worker records the last write or job it passed to Payssage. The 6-hourly `webhook.redrive` only runs if something happened since the last run, with a 15-minute margin over the sweep's 2-minute `redriveAfter`. On idle days, 4 cron wakes that each chained into IdentityX become 0. |
| Next | **Check on day 3.** On an idle day, expect IdentityX only for the 03:00 cron (about 5 minutes), Payssage and Informd at 0, and Univents only with real use. The target is under 0.3 GiB-h a day when idle, about 7 GiB-h a cycle, leaving the rest for real users. |

## Follow-up: day 3 (2026-11-05)

24 hours from Nov 4 15:00 to Nov 5 15:00 UTC, ending before the edge-cache and cold-start deploys at 15:52 but including the Nov 4 deploys around 21:11. PR #70 works: the quiet night (Nov 5 00:00 to 13:00) cost about 0.04 GiB-h, well under the 0.3 target. The full day came to about 3.4 Durable Object hours (5.8 on day 1) and 2.7k GiB-s (about 0.75 GiB-h, vs. 3.1k on day 1), almost all from deploys and real use.

| App | DO awake | Container awake | Quiet night | Why |
|---|---|---|---|---|
| univents | 110 min | 98 min | 02, 04, 05 UTC (~2 min each) | Real use 13:00–14:00 (~60 min), the Nov 4 deploy, unexplained night wakes |
| identityx | 60 min | 55 min | 02, 03, 04 UTC (~2 min each) | The Nov 4 deploy (18:00–22:00), the 03:00 cron, waking along with Univents |
| payssage | 26 min | 16 min | 0 | The Nov 4 deploy only. Redrive skipped at 00, 06 and 12 UTC |
| informd | 9 min | 9 min | 0 | The Nov 4 deploy only |

- **The JWKS copy holds:** 3 `jwks:` keys in KV, and other APIs' cold starts no longer wake IdentityX.
- **The redrive skip holds:** `payssage:lastRedrive` is still the first run (around Nov 4 00:00), and idle crons since then didn't wake Payssage.
- **Workers:** 0 errors on every script. The API Workers saw about 1.2k requests, the fronts and portal about 3.3k.
- **Spend guard:** 11.6k GiB-s this cycle, about 3.2 of the 22.5 GiB-h cut-off, resetting Nov 24.

| Status | Change |
|---|---|
| Done | **Day-3 target met for idle hours.** Payssage and Informd stay asleep when idle. The only idle cost left is about 6 minutes of IdentityX and 6 of Univents a night. |
| Next | **Find what wakes Univents at night** (02, 04, 05 UTC) in the Worker logs, whether it's a cron, a job or a scanner, and whether it drags IdentityX awake at 02 and 04. |
| Next | **Rerun the queries after the edge-cache and cold-start deploys** (Nov 5 15:52–16:09), for Nov 6 15:00 to Nov 7 15:00, an idle window with no deploys. Cached GETs should answer without waking a container, with fewer failed cold starts. |
| Next | **Compare cold-start success** before and after the fixes, from the new tracking. |
| Next | **Watch headroom weekly until Nov 24.** At day 3's 0.75 GiB-h a day, the cycle ends near 17.5 of 22.5 GiB-h. Deploy days cost the most. |

## Follow-up: day 7 (2026-11-09)

24 hours from Nov 6 15:00 to Nov 7 15:00 UTC, an idle window with no deploys (the last was Nov 5 21:23). The four containers were awake about 20 Durable Object minutes in total (about 0.33 hours, vs. 3.4 on day 3), about 0.08 GiB-h. Every wake now ends in about 2 minutes. From Nov 6 to Nov 9, about three quarters of the awake time was scanners probing paths `isNoise()` doesn't catch.

| App | DO awake | Wakes | Why |
|---|---|---|---|
| identityx | 8.8 min | 3 | The 03:00 cron, real use at 20 UTC, Censys probing `/mcp` and `/jsonrpc` at 16 UTC |
| univents | 4.5 min | 2 | Real use at 20 UTC. The 02 UTC wake was a scanner (`/_environment`, `/_profiler/phpinfo`, `/info`) |
| payssage | 4.2 min | 2 | Both from a Metabase CVE scanner on `/api/session/properties`. Redrive still skipped |
| informd | 2.2 min | 1 | A scanner (`/_environment`, `/_profiler/phpinfo`) |

- **Univents' night wakes are scanners,** not a cron or a job, and they no longer drag IdentityX awake. Since the JWKS copy, no IdentityX wake lines up with them.
- **What still wakes containers:** probes without extensions or with encoded paths (`/user`, `/api`, `/sign-in`, `/graphql`, `/Dockerfile`, `/%252eenv`, `/%62ackend`). The WAF answers `.env`, `.git` and PHP paths with 403 at the edge, but these get past it and reach a container for a 404.
- **Cold starts:** 24 of 86 failed before the Nov 5 fixes (28%), 3 of 69 after (4%), and 2 of those 3 were in the deploy hour. Time to first answer dropped from about 4.9 s to 3.0 s. The process's own `boot` span is 0.7 to 1.0 s on most starts, so most of the rest is Cloudflare starting the container.
- **Guard vs. Cloudflare:** since Nov 2 12:12, the guard counted 13.3k GiB-s against Cloudflare's 12.4k GiB-s of reserved memory. It reads about 7% high, which is the safe direction. IdentityX matches within 1%, and the others read 11 to 22% high because the library reports some stops late.
- **Headroom:** the guard is at about 3.7 of the 22.5 GiB-h cut-off. At about 0.1 GiB-h per idle day, the cycle ends near 5 to 7 GiB-h on Nov 24, with deploy days the main cost.

| Status | Change |
|---|---|
| Done | **Day-3 checks closed.** Idle cost is down about 10× from day 3, failed cold starts about 7×, and the guard errs on the safe side. The incident is contained. |
| Next | **Only forward real routes.** Replace the `isNoise()` denylist with an allowlist: each API Worker forwards a path only when its first segment is a top-level path in that API's spec (like `/projects`, `/auth` or `/.well-known` for IdentityX), and answers 404 otherwise. Every scanner wake above fails that check and all real traffic passes, leaving the idle cost at the 03:00 cron plus real use. |

*Written 2026-11-02, with follow-ups on 2026-11-03, 11-05 and 11-09, from Cloudflare billable usage, GraphQL Analytics, the Containers API, Analytics Engine (`container_cold_starts`), Grafana Tempo (`boot` spans), the `@cloudflare/containers` source and git history (PRs #68, #69 and #70).*
