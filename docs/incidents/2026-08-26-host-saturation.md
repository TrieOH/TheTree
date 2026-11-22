---
title: A publish build froze the trieoh host
date: 2026-08-26
status: resolved
severity: high
started: 2026-08-26T17:04Z
ended: 2026-08-26T17:34Z
summary: A normal CI build pinned the 2-core VPS at 100% CPU for 18 minutes, and the runner's retries made Forgejo balloon until the build finished. No data was lost.
services: [forgejo, forgejo-runner, all hosted services]
---

## What happened

A routine publish build (`univents/v0.10.12`, commit `d9dec307`, tagged at 17:03 UTC) ran the Dagger `publish` workflow inside `forgejo-dind` on the trieoh host. The host had 2 cores, 7.8 GiB of RAM and no swap. dind was capped at 1.5 CPU and 3 GiB.

The build succeeded, but for about 18 minutes it held the host at 100% CPU. Forgejo slowed down so much that the Actions runner couldn't report for 15 minutes, and every service stopped responding. When the build finished at 17:34, the host recovered by itself.

In short: a CI job took 75% of the host's two cores on a machine already running everything else, and the runner's 1-second reporting turned Forgejo's slowness into a retry storm. Forgejo grew to about 2 GB and 90% CPU, which made it slower, which caused more retries, until the build ended.

## Timeline

All times UTC, minute by minute from beszel.

| Time | Host CPU | Host RAM used | What happened |
|---|---|---|---|
| 16:45–17:04 | 30–50% | 2.5 G (32%) | Normal, load around 1 |
| 17:04:24 | | | Runner picks up task 1213 (`TrieOH/TheTree`, publish workflow) |
| 17:05–17:06 | 90 → 100% | 3.0 → 3.5 G | 1-minute load jumps to 17.7, disk starts growing |
| 17:06–17:23 | 98–100% | up to 5.08 G (65.5%) | Forgejo grows to 1.9–2.0 GB and ~90% CPU. Runner reports time out from 17:07 to 17:22. Load peaks at 21.7. Disk goes from 31 to 37.7 G |
| 17:24–17:29 | 34–47% | 4.5 → 4.0 G | Heavy part ends, Forgejo back to ~570 MB, reports go through again |
| 17:30–17:33 | 84–93% | 4.3–5.1 G | Image export and push. dind hits its 3 GiB cap (3004 MB at 17:31) |
| 17:34 | 29% | 2.8 G (36%) | Build done, Dagger engine exits 0, host recovers |

## Numbers

- **CPU:** 100% for about 18 minutes on 2 cores, 1-minute load peaking at 21.7.
- **RAM:** real memory use peaked at 5.08 of 7.75 GiB (65.5%). It never ran out: no OOM kills, no container crashed.
- **Why it looked like 100% RAM:** about 2.6 GiB of page cache sat on top. With no swap and 8.6 GB written to disk in 30 minutes, the kernel kept reclaiming cache, so everything thrashed.
- **Disk:** 30.9 → 39.5 G during the build, 8.6 G of images and BuildKit cache.
- **Forgejo memory:** it grows from about 250 MB to over 550 MB each day until the 05:00 restart. It was already at 564 MB before the build.
- **Result:** the build succeeded, the image was pushed, nothing was lost.

## Why

1. **Trigger:** a normal publish build. Nothing ran away. It was just a build on a machine with no room to spare.
2. **Sharing the host:** everything ran on the same 2-core, 7.8 GiB, no-swap box as CI: Forgejo, Postgres, Grafana, three Victoria services, four apps, mox, ntfy and beszel, using 2.5 to 3 GiB at rest. dind's 1.5-core cap was 75% of the CPU.
3. **The retry storm:** the runner used `report_interval: 1s` with short timeouts. When Forgejo slowed, the runner kept retrying, Forgejo grew to about 2 GB and 90% CPU, got slower, and the retries got worse. Forgejo, not the build, used the most CPU during the incident.

## Evidence

- beszel history (SQLite, 30 days): host and per-container CPU and RAM per minute, the source of the timeline
- Runner logs: task 1213 picked up at 17:04:24, then `context deadline exceeded` reporting to `forgejo:3000` from 17:07 to 17:22
- Kernel logs: no OOM kills in the window
- `docker inspect`: the dind limits were applied (Compose v5.3.1 honors `deploy.resources.limits`)
- The Dagger engine container exited 0 at 17:34

## What we changed

| Status | Change |
|---|---|
| Done | **Stop the retry storm.** `forgejo/runner/config.yml`: `report_interval` 1s → 5s, `fetch_timeout` 5s → 20s, `fetch_interval` 2s → 5s. |
| Done | **Smaller dind.** `forgejo/compose.yml`: 1.5 CPU / 3 GiB → 1.0 CPU / 2.5 GiB, so Forgejo can always answer the runner. Builds are a bit slower and may sometimes hit the memory cap. |
| Done | **Cap Forgejo's memory** at 2 GiB in `forgejo/compose.yml`, so one bad day can't eat the host. The 05:00 restart stays. |
| Done | **Host alerts in Grafana.** `observability/` gained node-exporter and two alerts, CPU over 90% and memory over 85% for 5 minutes. The CPU alert would have fired around 17:11. Grafana still needs a contact point set in the UI to send them. |
| Done | **beszel alerts** for CPU over 90%, memory over 85% and disk over 85%, averaged over 5 minutes, emailed to admin@trieoh.com. ntfy push can be added in the beszel UI. |
| Done | **4 GiB of swap** at `/swapfile`, in `/etc/fstab` with priority -2, so the host has a buffer under load. The script is `infra/scripts/add-swap.sh`. |
| Done | **Forgejo's "leak" explained.** Forgejo 10.0.3 run locally under heavy load for an hour kept memory, Go heap and goroutines flat, so there's no heap leak. In production, about 220 MB of the 430 MB cgroup was page cache the kernel never had to drop, and the app itself stayed at about 191 MB. The daily growth (264 → 466 MB) stopped after the runner fix: the old 2-second fetch and 1-second report polling caused it. The daily restart stays as a backup. |
| Later | **Run Actions on another host.** Two cores running production and CI together will keep hitting this. |

*Written 2026-08-26 from beszel history, runner logs, kernel logs, `docker inspect` and git history.*
