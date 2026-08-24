# Infra Wishlist

Status: **open — nothing built yet** (except #46, done 2026-08-24). Date: 2026-08-24. Scope: `TrieOH/infra` · `TrieOH/deploy` · `TrieOH/TheTree` frontends · Cloudflare (Workers/DNS). Server: `trieoh@main`, single VPS (8GB RAM), Docker Compose, no Swarm.

Everything here is single-server-shaped: cheap, compose-native, and mostly "wire existing pieces together". Items are grouped **by owner** — who must do the work — not by theme. Every item keeps its number (1–39) so references stay stable; the **order of operations** at the end is the sequencing guide regardless of owner.

**The split in one line:** ~30 items are pure infra/devops (config, compose, host, DNS, CI tooling). Two items need a developer (#11 error tracking, #16 e2e). Two are shared (#9 log shipping, #18 staging telemetry) — and **#9's dev half is the prerequisite for #11 and #18**, so schedule it with the infra wiring. Plus **5 infra-DX items (#35–39)** — rebuild/move ergonomics — and **6 infra simplifications (#40–45)** — deletion over addition.

---

## 🧰 Infra/DevOps — do it solo (30 items)

| # | Item | Effort |
|---|------|--------|
| 1 | Offsite automated backups (restic → R2/B2) | M |
| 2 | Grafana alerts → ntfy (phone) | S |
| 3 | External uptime checks as a Cloudflare Worker cron (+ status page) | S |
| 4 | Container log rotation / disk-full alert | S |
| 5 | zram / swapfile | S |
| 6 | Profile container memory before setting limits | XS |
| 7 | Resource limits on every compose stack | S |
| 8 | TLS certificate expiry alert | S |
| 10 | Latency alerts on existing traces | S |
| 12 | Auth gate on staging (CF Access + Caddy basic_auth) | S |
| 13 | Wildcard DNS for staging (4 records, not 8) | XS |
| 14 | mailpit retention flags | XS |
| 15 | Nightly staging reset (Forgejo schedule) | S |
| 17 | Compose drift check prod ↔ staging | S |
| 19 | Server crontab into git (`cron/` + `just install-crons`) | S |
| 20 | Host image GC (beyond DinD) | XS |
| 21 | unattended-upgrades (security-only) | S |
| 22 | SSH hardening (key-only, fail2ban) | S |
| 23 | Host firewall (ufw, default-deny) | S |
| 24 | Secret scanning (gitleaks) — CI tooling | S |
| 25 | DNS-as-code via the CF API | S |
| 26 | Dockge (compose web UI) | M |
| 27 | Auto-bump PRs in the deploy repo | M |
| 28 | `docker compose up -d --wait` in deploy recipe | XS |
| 29 | Validate the Caddyfile before reload | XS |
| 30 | Server ops kit in the infra justfile (+ rollback paths) | S |
| 31 | Pin `informd` | XS |
| 32 | Tailscale (replace SSH tunnels) | M |
| 33 | PgBouncer — only if connection errors appear | XS |
| 34 | PITR via WAL archiving | L |

## 🤝 Shared — infra + a developer (2 items)

| # | Item | Effort | Split |
|---|------|--------|-------|
| 9 | Log shipping → VictoriaLogs | M | **dev:** OTLP logs exporter in `lib/go/telemetry` · **infra:** endpoint/datasource, Caddy access-log routing, postgres slow-query config |
| 18 | Staging telemetry → same Victoria instances, `env` labels | S | **infra:** obs-net wiring, `ENV=staging`, alert filters · **dev:** `env` label in the Go telemetry config |

> **Dependency:** #9's dev half (the logs exporter) gates #11 and #18. The Go change and the infra wiring must land in the same window — otherwise the whole observability chain stalls.

## 👨‍💻 Developer — not infra work (2 items)

| # | Item | Effort | Why it's dev |
|---|------|--------|--------------|
| 11 | Error tracking = error logs + trace error spans | S | `RecordError` helper in `lib/go/telemetry` + calls at error sites across four services + ServiceVersion wiring — application code. (Infra slice: the Grafana panel + alert rule.) |
| 16 | E2e package (Playwright) against staging | M | Playwright tests against staging are test/application work — no infra component. |

## 🧰 Infra DX — rebuild, move & ergonomics (5 items, all infra-owned)

The things that make the box *movable*: rebuild from scratch, relocate to another server, or just operate without a checklist in someone's head. These pair with existing items — #35 with #1, #36 with #5/#19/#21/#22/#23, #37/#39 with each other.

| # | Item | Effort | Why it matters for DX |
|---|------|--------|------------------------|
| 35 | Off-box secrets vault (age/sops) | M | Today every secret lives only on the box — a move without the old box is impossible |
| 36 | `just bootstrap` — whole box from nothing | M–L | setup.sh covers only caddy+forgejo; obs-net is created by no script |
| 37 | Rebuild runbook doc | S | The ordering (registry chicken-and-egg, DNS first) is written nowhere |
| 38 | Cert provisioning into git | S | Zero docs on how the wildcard + mox certs are obtained/renewed |
| 39 | Rebuild drill (annual) | M | A runbook nobody has executed is a fiction |

## 🧰 Infra simplifications — deletion over addition (6 items, all infra-owned)

| # | Item | Effort | What it deletes |
|---|------|--------|-----------------|
| 40 | Delete dead Caddyfile config | XS | `api.trieoh.com` block + `(api_routes)` snippet — zero consumers (verified) + mailpit comment |
| 41 | Snippet-ify the repeated mail blocks | S | 3× acme-challenge handlers + 3× identical letsencrypt tls lines |
| 42 | Consolidate cert provisioning to one source | M | The `/etc/letsencrypt` path + hand-rolled acme-challenge blocks vanish |
| 43 | Kill the forgejo-restart cron by fixing its root cause | M | A nightly restart is a symptom patch; delete the cron |
| 44 | Rename forgejo's confusing `internal` network | XS | Same name as thetree's `internal`, different meaning (`internal: false`!) |
| 45 | `just register-runner <token>` script | XS | The manual `docker run` registration dance in the README |
| 46 | Forgejo composite action — pinned, cached docker CLI | S | 5× copy-pasted `curl \| sh` installs a redundant daemon, unpinned, downloads every run |

---

# Details — 🧰 Infra/DevOps (solo)

## Protect the box

### 1. Offsite automated backups

**Problem.** Nothing automated backs up the box. The only mentions are "keep a baseline in `~/backups/`" (deploy README, ADR-0007) — a manual folder on the same disk it's protecting. One disk failure loses postgres, rustfs, and the env files all at once.

**Fix.**
- Nightly `pg_dump` of the postgres volume (or `pg_basebackup` for a full-cluster copy), plus a restic snapshot of `rustfs_data` / `rustfs_logs`.
- restic → **Cloudflare R2** (account + creds pattern already exists for workers) or Backblaze B2. Encrypted by restic.
- Push a **ntfy failure notification** when a backup errors — a silent backup is a non-backup.
- Also snapshot the 6 env files (`.env` + `.identityx.env` / `.univents.env` / `.payssage.env` / `.informd.env`) **age/sops-encrypted** to the same target. They're server-only by design (ADR-0007); losing the box loses the secrets.
- **Scope is explicit:** include postgres, rustfs, and the env files; **exclude** the Victoria volumes — metrics/logs/traces are 7d-ephemeral and rebuildable. Nobody should assume "everything" is backed up.
- Retention: e.g. 7 daily / 4 weekly / 6 monthly, prune nightly.
- Do one **restore drill** — an untested backup is a hope.

**Where.** New `infra/backup/` (script + systemd timer or crontab entry — see #19) + R2 bucket/API key.

**Effort.** M. Dependencies: R2 or B2 creds; restic container or binary.

### 2. Grafana alerts → ntfy

**Problem.** `infra/observability/provisioning/alerting.yml` provisions rules (service-down, high-5xx, obs-stack-down) but **no contact points / notification policies** — alerts fire inside Grafana and wake nobody. `ntfy.trieoh.com` already pushes to your phone; the two are not wired.

**Fix.** Add a provisioned contact point (Grafana → ntfy webhook, or email via mox) + notification policy routing `severity=critical|warning` to it. Provisioning means it survives a Grafana volume wipe — don't configure it by clicking in the UI.

**Where.** `infra/observability/provisioning/alerting.yml` (+ maybe `contact-points.yml`).

**Effort.** S.

### 3. External uptime checks — Cloudflare Worker cron

**Problem.** The original idea was a server cron curling `/health`. A check living on the box dies with the box — the worst moment to go blind. Grafana is on the same box too.

**Fix.** A tiny **Cloudflare Worker with a Cron Trigger** (`triggers: [{ cron: "*/5 * * * *" }]`), running globally — outside the box. It pings each public `/health` endpoint, and only when one fails, POSTs to ntfy. Checks the *public surface* (Caddy + service), which is what users experience. Later the staging endpoints (`staging.api.*`) just get added to the same list.

Sketch (`uptime-check` worker, e.g. `infra/uptime-check/`):

```ts
const ENDPOINTS = [
  { name: "identityx", url: "https://api.identityx.com.br/health" },
  { name: "univents", url: "https://api.univents.com.br/health" },
  { name: "payssage", url: "https://api.payssage.trieoh.com/health" },
  { name: "informd", url: "https://api.informd.trieoh.com/health" },
  { name: "forgejo", url: "https://git.trieoh.com/api/healthz" },
];

export default {
  async scheduled(_event, env, _ctx) {
    const failures = [];
    for (const ep of ENDPOINTS) {
      try {
        const res = await fetch(ep.url, { signal: AbortSignal.timeout(10_000) });
        if (!res.ok) failures.push(`${ep.name}: HTTP ${res.status}`);
      } catch (e) {
        failures.push(`${ep.name}: ${(e as Error).message}`);
      }
    }
    if (failures.length > 0) {
      await fetch(env.NTFY_URL, {
        method: "POST",
        headers: { Title: "Uptime alert" },
        body: failures.join("\n"),
      });
    }
  },
};
```

Wrangler: `"triggers": [{ "cron": "*/5 * * * *" }]`; secret `NTFY_URL = https://ntfy.trieoh.com/<topic>` (add a dedicated topic + optional access token in `infra/ntfy/config/server.yml`).

**Bonus — status page:** the same worker can serve `https://status.trieoh.com` — a green/red page rendered from the last run's results (a plain `fetch` handler on the default route). Zero extra infra, gives users a "we're up" page without adding Gatus or anything else to the box.

**Where.** New `infra/uptime-check/`. Deploy with the existing wrangler-action pattern.

**Effort.** S.

### 4. Container log rotation + disk-full alert

**Problem.** No `logging:` config anywhere → Docker's default `json-file` driver grows unbounded. Disk-full is the classic single-VPS death (postgres stops writing, everything cascades). And since logs don't flow to VictoriaLogs (see #9), the json-file copies are all you have — but they still need capping.

**Fix.** Per-service `logging: options: { max-size: "50m", max-file: "3" }` in every compose file (once #9 lands, backends can even go `driver: none`). Plus a **disk usage alert**: Beszel already monitors the host — wire its disk alert to ntfy (Beszel supports ntfy/webhook out of the box).

**Where.** `deploy/thetree/compose.yml` (+ staging), `infra/*/compose.yml`, `infra/beszel/`.

**Effort.** S.

### 5. zram / swapfile

**Problem.** 8GB box running prod + CI (DinD) with staging to come — memory spikes are a when, not an if. With no swap, an overcommit OOM kills whichever process the kernel picks — historically the one holding the most pages, which is often postgres or Caddy. The two things that must never die.

**Fix.** Check first: `swapon --show`. If empty, add **zram** (compressed RAM swap, ~1–2GB, near-zero latency cost) or a plain 2GB swapfile. The kernel then thrashes instead of killing the wrong process, and `restart: unless-stopped` recovers the victim container. Optional upgrade: `systemd-oomd` for smarter kill selection.

**Where.** Host config, documented in `infra/` (setup script + rebuild checklist, #19).

**Effort.** S. One-time.

### 6. Profile container memory before setting limits

**Problem.** #7's limits need real numbers, not guesses — "512m for backends" is a guess until measured.

**Fix.** Before staging lands: `docker stats --no-stream` snapshot (or a 60s `docker stats` sample) per container, record peak RSS in the deploy repo (a `MEMORY.md` or comments in compose). Re-measure once staging is up. Then set #7's limits from the data — including headroom for CI spikes (DinD builds).

**Where.** One-time record, `deploy/thetree/` or `infra/`.

**Effort.** XS.

### 7. Resource limits on every compose stack

**Problem.** The box runs 6 compose stacks (caddy, forgejo, mox, ntfy, beszel, observability) + the thetree prod stack, with the staging twin to come — and **only `forgejo-dind` has any limits** (`cpus: "1.5"`, `memory: 3g`). Every other container is unbounded. One runaway process (a CI build, a leaky service, a Grafana query storm, a staging spike) can OOM the host and take everything down — including postgres and Caddy, the two things that must never die.

**Fix.** Do **#6 first** (measure), then `deploy.resources.limits` on every service in every compose file, using the `forgejo-dind` pattern:

```yaml
deploy:
  resources:
    limits:
      cpus: "0.5"    # per-service budget, tuned per workload
      memory: 512m
```

Suggested starting budgets (adjust from #6's measurements): backends 512m–1g each, postgres 1–2g, rustfs 512m, caddy/mox/ntfy/beszel 256m each, observability 512m each, staging services *tight* (256m) so staging can never starve prod. The failure mode becomes "one container gets OOM-killed and `restart: unless-stopped` brings it back" instead of "the box dies". Set `memory_reservation` (soft) too, so the kernel schedules fairly before anyone hits a hard limit.

**Where.** `deploy/thetree/compose.yml` (+ staging) + every `infra/*/compose.yml`.

**Effort.** S. The cheapest insurance on the list — a config change, no new moving parts.

### 8. TLS certificate expiry alert

**Problem.** The box has **three cert sources**: Caddy's automatic Let's Encrypt, `/etc/letsencrypt` (mox), and `/etc/caddy/certs` (the `*.trieoh.com` wildcard). Each renews on a different mechanism. If certbot or the wildcard renewal breaks, sites go dark with **zero logs and zero alarms** — the classic "worked for a year, silently died on a Sunday".

**Fix.** A cron (joins the #19 crontab) that checks every public host's cert end-date and ntfy-alerts at <14 days:

```sh
for h in api.identityx.com.br api.univents.com.br api.payssage.trieoh.com \
         api.informd.trieoh.com git.trieoh.com ntfy.trieoh.com \
         mail.trieoh.com mta-sts.trieoh.com; do
  end=$(echo | openssl s_client -servername "$h" -connect "$h":443 2>/dev/null \
        | openssl x509 -noout -enddate | cut -d= -f2)
  # ...diff against today, alert to ntfy if < 14 days...
done
```

Note: **mox's certs are a separate renewal path** (`/etc/letsencrypt`) — the list must include `mail.trieoh.com` / `mta-sts.trieoh.com`, not just the gateway hosts. (Alternative: fold a cert check into the #3 uptime worker — but `fetch()` in workers doesn't expose cert dates, so the cron is the honest place for this one.)

**Where.** `infra/cron/` + ntfy topic.

**Effort.** S.

## Observability

### 10. Latency alerts on the traces you already have

**Problem.** Traces flow end-to-end today (otelhttp wraps every router in `lib/go/httpserver`, backends ship OTLP to victoria-traces, frontends too). That investment is invisible — dashboards exist but there's no latency alert.

**Fix.** A p99-per-route alert in `alerting.yml` (e.g. `histogram_quantile(0.99, ...)` over the otel http.server duration metric, threshold per service, `for: 10m`), routed to the same ntfy contact point as #2. Low effort, big signal.

**Where.** `infra/observability/provisioning/alerting.yml`.

**Effort.** S.

## Staging

### 12. Auth gate on staging

**Problem.** Per the staging plan, `staging.identityx.com.br` etc. go live on the public internet with real subdomains — anyone can sign up, create actors, trigger emails, hammer the API. The plan has no auth gate.

**Fix.**
- **Frontends** (CF workers, custom domains): **Cloudflare Access** (free tier = 50 users). A config toggle on the account, no code — anyone without a `@trieoh.com` (or whatever) email gets a login wall.
- **API subdomains** (grey-cloud → Caddy): `basic_auth` on the 8 staging blocks — the `(admin_auth)` snippet pattern already exists in the Caddyfile.
- E2e tests that need to pass the gate use a service token / bypass policy scoped to their IP or a header.

**Where.** Cloudflare dashboard + `infra/caddy/Caddyfile` (staging blocks).

**Effort.** S. Do this the same day the staging DNS records go in.

### 13. Wildcard DNS for staging

**Problem.** The plan lists 8 A records (`staging.api.*` + `staging.auth.*` × 4 zones). Every future staging service adds more.

**Fix.** One `*.staging.<zone>` A record per zone instead — `*.staging.identityx.com.br`, `*.staging.univents.com.br`, `*.staging.payssage.trieoh.com`, `*.staging.informd.trieoh.com` — all grey-cloud (proxied off, so Let's Encrypt works). Covers `api`, `auth`, and anything added later, with the same 4 records.

**Where.** DNS (4 records, not 8).

**Effort.** XS.

### 14. mailpit retention flags

**Problem.** Staging mailpit runs forever (only reset on `--reset`); its DB grows with every e2e run's emails.

**Fix.** Cap it in the staging compose — mailpit flags: `--max-messages 500 --max-age 7d` (approx; verify flag names for v1.30.x).

**Where.** `deploy/thetree/staging/compose.yml`.

**Effort.** XS.

### 15. Nightly staging reset

**Problem.** Staging data accumulates; e2e assertions silently rot against state nobody remembers creating.

**Fix.** A **Forgejo Actions `schedule:` workflow** (Forgejo supports cron schedules) that SSHes to the box and runs `cd ~/deploy/thetree/staging && ./bootstrap.sh --reset` nightly. e2e always starts from a known state. (Alternative: the #3 CF worker gains a second cron that hits a reset endpoint.)

**Where.** New workflow in `TrieOH/deploy` (or TheTree) + a server SSH key for the runner.

**Effort.** S.

### 17. Compose drift check prod ↔ staging

**Problem.** Two full compose files (prod + staging) diverge by hand — add a service to prod, forget staging, and e2e silently tests an outdated twin.

**Fix.** A one-liner in the deploy repo CI (or a pre-commit hook): normalize both via `docker compose config`, strip the *intended* deltas (project name, container names, ports, volumes, networks), and diff the rest — so a diff only surfaces real drift (missing services, different image tags, env shape).

**Where.** `TrieOH/deploy` (workflow or hook).

**Effort.** S.

## Server hygiene

### 19. Server crontab into git

**Problem.** The deploy README: "The nightly dind-prune and forgejo-restart crons live in the server's crontab (not in git) — when rebuilding the box, re-add them." That's rebuild-debt: the box is not reproducible from the repo.

**Fix.** A `cron/` directory in the infra repo with the entries as files, plus `just install-crons` that installs them (idempotently) into the server crontab. After #3, the uptime check moves off the box entirely; what remains is dind-prune, forgejo-restart, the #8 cert check, and (after #1) the backup job. Pair with a short rebuild checklist (crontab, SSH keys, letsencrypt dirs, swap from #5).

**Where.** `infra/cron/` + `infra/justfile`.

**Effort.** S.

### 20. Host image GC (beyond DinD)

**Problem.** The nightly prune covers the DinD cache; the host itself accumulates every pulled registry tag (`identityx:v0.35.3`, old builds, staging twins) on disk.

**Fix.** Add `docker image prune -a --filter until=168h` to the same cron as #19.

**Where.** `infra/cron/`.

**Effort.** XS.

### 21. unattended-upgrades (security-only)

**Problem.** A lone VPS with open SSH (see #22) is only as patched as you remember to patch it.

**Fix.** `unattended-upgrades` restricted to security updates; kernel reboots stay manual (or a monthly reboot cron). Tradeoff: surprise package upgrades — keep it security-only to bound the risk.

**Where.** Host config (`infra/` docs or a setup script).

**Effort.** S. Optional but standard.

### 22. SSH hardening

**Problem.** `trieoh@main` exposes plain SSH (port 22) to the internet — the one public surface beside Caddy.

**Fix.** Key-only auth (password auth off) + fail2ban/sshguard for brute force. If #32 (Tailscale) lands, the stronger move is restricting SSH to the tailnet entirely — then hardening is mostly moot.

**Where.** Host config (`/etc/ssh/sshd_config` + fail2ban), documented in `infra/`.

**Effort.** S.

### 23. Host firewall (ufw, default-deny)

**Problem.** Whatever compose maps a port to is exposed: today that's 22, 80, 443 + mox's 25/465/587/993. No default-deny — a stray service binding 0.0.0.0 (or a future compose mistake) is immediately public.

**Fix.** `ufw default deny incoming` + allow exactly: 22 (SSH), 80/443 (Caddy), 25/465/587/993 (mox SMTP/IMAP). One-time, standard, and it bounds the blast radius of any future port mistake. Also review after each new compose service.

**Where.** Host config, documented in `infra/` (setup script + rebuild checklist).

**Effort.** S.

### 24. Secret scanning (gitleaks)

**Problem.** ADR-0007 exists because a `docker compose config` snapshot once leaked **every prod secret into git** — scrubbed after the fact with `git filter-repo` + force-push. Today the defense is a human rule ("never commit snapshots"). Humans lapse; the deploy repo is one bad paste away from a repeat.

**Fix.** **gitleaks** (free, single binary) — CI tooling, devops scope even though the workflows live in the app repo:
- CI on push/PR in **`TrieOH/deploy`** (where the templates live and the incident happened) and **TheTree** — fail the build on any secret-shaped string.
- Pre-commit hook locally (TheTree already has Husky — add a gitleaks stage).
- One pass over history (`gitleaks detect --log-opts=--all`) to confirm the old scrub is actually clean.
- Tune an `allowlist` for test tokens/placeholders so it doesn't noise-cry.

**Where.** Workflows in `TrieOH/deploy` + TheTree; `.husky/pre-commit` in TheTree.

**Effort.** S. Directly prevents a recurrence of your one real infra incident.

### 25. DNS-as-code via the CF API

**Problem.** DNS exists only in the Cloudflare dashboard, edited by hand. Three zones today, staging records coming (#13) — churn is growing, and the config isn't reproducible.

**Fix.** A `records.yml` (or `dns/*.yml`) in the infra repo + `just dns-sync` — the `cloudflare` CLI or a small script against the CF API applying the declared records. DNS becomes reviewable and rebuildable like everything else. Optional at this scale, but cheap once the staging records land.

**Where.** `infra/` + a CF API token (DNS edit scope).

**Effort.** S.

### 26. Dockge (compose web UI)

**Problem.** Every ops action is a terminal + SSH session. Sometimes (phone, quick check) a UI beats a terminal.

**Fix.** **Dockge** — a compose-native web UI: starts/stops/restarts stacks, tails logs, edits compose files, all stored as plain compose files (no proprietary state). Fits the repo layout (one folder per stack). Add basic auth (or Tailscale-gate it) since it's a web surface with compose edit rights.

**Where.** New `infra/dockge/` (or under Tailscale only).

**Effort.** M. Low priority — nice-to-have, not a gap.

## Deploy ergonomics

### 27. Auto-bump PRs in the deploy repo

**Problem.** Bumping a version = manually editing `deploy/thetree/compose.yml` (the "git history is the version ledger" flow). Fine, but it's the most common deploy action, done by hand.

**Fix.** After `publish.yml` pushes a tag, a workflow opens a **PR in `TrieOH/deploy` bumping the pin** (never auto-merge — you still review, but the edit is generated). Needs a Forgejo token with repo-write. Cheap version: an ntfy ping "identityx v0.36.0 published" and you bump by hand.

**Where.** Workflow in TheTree (on publish) or deploy (polling the registry).

**Effort.** M. Do the ntfy-ping version first if the token wiring feels heavy.

### 28. Healthcheck-gated deploys

**Problem.** The deploy recipe is `docker compose up -d --no-deps` — compose returns once containers start, not once they're healthy.

**Fix.** `docker compose up -d --wait` (compose's `--wait` already waits on the declared healthchecks — all services have them). Fail loudly if a service never becomes healthy instead of discovering it on the next curl.

**Where.** `deploy/README.md` (recipe), future deploy workflow.

**Effort.** XS.

### 29. Validate the Caddyfile before reload

**Problem.** `just reload-caddy` applies whatever's on disk (`docker exec caddy caddy reload --config /etc/caddy/Caddyfile`). A malformed Caddyfile (typo, stray brace, broken snippet) takes **the whole gateway down — every site at once**. This failure class already bit once (the observability compose comment: "latest" drift broke Caddy→Grafana DNS → 502).

**Fix.** Make reload refuse broken config:
- `just reload-caddy` runs `caddy validate --config /etc/caddy/Caddyfile` first and exits non-zero on failure — nothing gets applied.
- `caddy fmt --check` (or a validate step) as a pre-commit hook / CI check on the infra repo, so a broken Caddyfile never reaches the box.

**Where.** `infra/justfile`, infra pre-commit or workflow.

**Effort.** XS.

### 30. Server ops kit in the infra justfile (+ rollback paths)

**Problem.** "How do I look at prod" is tribal knowledge: which commands, which dirs, which tunnels. And when something breaks, the rollback paths aren't written down anywhere central.

**Fix.** Extend the infra justfile with recipes: `logs <svc>`, `ps`, `status` (healthchecks), `backup now`, `restore-drill`, plus the rebuild checklist from #19. Also document the two rollback paths in the same surface:
- **Frontends** — Cloudflare dashboard → Workers → Versions: one-click rollback to the previous deployment.
- **Backends** — `git checkout <prev> -- compose.yml` + redeploy (already in the deploy README; make it a `just rollback-backend <tag>` recipe).

**Where.** `infra/justfile` + `infra/README.md`.

**Effort.** S.

### 31. Pin `informd`

**Problem.** `deploy/thetree/compose.yml` runs `informd:latest` — an explicit escape hatch (deploy README: "pin it on its next release"). Unpinned = silent drift on redeploy.

**Fix.** On informd's next release, pin the tag like the other three (`identityx:v0.35.3` pattern).

**Where.** `deploy/thetree/compose.yml` (+ staging mirror).

**Effort.** XS.

## Access

### 32. Tailscale

**Problem.** Admin access today is SSH + ad-hoc tunnels (`ssh -N -L 8026:127.0.0.1:8026` for mailpit, loopback debug ports for scripts). Each tunnel is a manual command; debug ports are only on loopback because there's no safe path to them otherwise.

**Fix.** **Tailscale** on the box + your machines:
- Replace SSH tunnels with a permanent, encrypted path: mailpit UI, staging debug ports (`127.0.0.1:18080/18082` per the plan), and admin endpoints are reachable by Tailnet IP/name only — not exposed to the internet at all.
- Tailscale SSH can replace plain SSH (key policy managed in one place) — which also solves #22.
- Free for a single user; no open ports added to the box.
- Keep SSH as-is if you prefer — the tunnel-removal is the win, not the migration.

**Where.** `infra/` (a `tailscale` compose service or the host package; the agent runs on `trieoh@main`).

**Effort.** M (mostly setup + deciding what stays on loopback).

## Postgres

### 33. PgBouncer — only if you see it

**Problem.** Eight services (prod + staging) will share one PG instance. If you ever hit "too many connections" errors, that's the signal.

**Fix.** Add PgBouncer in front of postgres (a compose service, pointed at by the 8 backends). Don't preemptively add it — Go's `database/sql` pooling is fine until it isn't, and a pooler is one more hop to debug.

**Where.** `deploy/thetree/` (+ staging).

**Effort.** XS until needed, then S. Deferred by design.

### 34. PITR via WAL archiving

**Problem.** Nightly dumps (#1) lose up to 24h of data. For a SaaS with payments (Payssage/MercadoPago), that window may or may not be acceptable.

**Fix.** `pg_basebackup` + WAL archiving (to R2/B2) enables point-in-time recovery. This is the *stretch* goal: more moving parts, restore is more involved. Do #1 and the restore drill first; revisit PITR only if the dump window actually hurts.

**Where.** `infra/backup/`.

**Effort.** L. Deferred unless the RPO bites.

---

# Details — 🤝 Shared (infra + a developer)

### 9. Log shipping → VictoriaLogs

**Problem.** **Confirmed gap:** `lib/go/telemetry` exports OTLP **traces** only — there is no log exporter anywhere in the four backends. VictoriaLogs and its Grafana datasource plugin are running, but nothing feeds them. Your zap logs only land in Docker's json-file: unsearchable, unalertable, and the reason #4 matters.

**Who.**
- **Developer:** the **OTLP logs exporter** in `lib/go/telemetry` — Go SDK code (mirrors the existing `InitTracer` pattern, endpoint `http://victoria-logs:9428`).
- **Infra:** the pipeline around it — datasource wiring, Caddy access-log routing, postgres slow-query config, stream labels.
- **⚠️ This dev half is the prerequisite for #11 and #18** — land the Go change in the same window as the infra wiring.

**Fix.** Two options:
- **OTLP logs exporter in `lib/go/telemetry`** — one place, all four services adopt it. Cleanest.
- **Vector sidecar** per service — no code change, but 4 more containers and config duplication.

**More streams once the pipeline works (infra):**
- **Caddy access logs** — the `log` directives in the Caddyfile vanish into docker json-file today; 5xx triage ("is it Caddy or the service?") is invisible. This already bit once (the "latest" drift broke Caddy→Grafana DNS → 502, observability compose comment). Route Caddy's logs into the same pipeline with the same `env` labels.
- **Postgres slow queries** — `log_min_duration_statement = 500ms` → stderr → the same pipeline. "Why is this endpoint slow" becomes a LogsQL query instead of guesswork.

Payoff: searchable logs in Grafana, log-based alerts, 7d retention instead of unbounded disk — and the foundation for #11 (error tracking with zero new infra).

**Where.** `lib/go/telemetry` (dev) + `infra/observability` / `infra/` / `deploy/thetree/` (infra).

**Effort.** M. Highest-value observability item.

### 18. Staging telemetry → same Victoria instances, `env` labels

**Problem.** The staging plan's "no `obs-net`" delta (staging must not pollute prod observability) has a hidden cost: **staging becomes blind**. No logs, no metrics, no traces — so #9/#10/#11 (log shipping, latency alerts, error tracking) only ever cover prod, and e2e failures against staging are undebuggable. The alternative — a second Victoria stack — doesn't fit 8GB.

**Who.**
- **Infra:** `*-staging` on `obs-net`, `ENV=staging` env, alert filters (`env="prod"` for prod rules).
- **Developer:** the `env` label handling in the Go telemetry config.

**Fix.** Put `*-staging` on `obs-net` and give every stream an **`env` label**: VictoriaLogs `_stream:{env="staging", service="identityx"}`, same for metrics/traces. Then filter all prod alerts to `env="prod"` (and staging alerts — or none — separately). Same instances (that's the RAM win), clean separation, and staging errors are visible in the same Grafana. This is the difference between "staging is a twin" and "staging is a black box". *(Amends the staging plan — note added to `docs/plans/staging-environment.html`.)*

**Where.** `deploy/thetree/staging/compose.yml` (infra), telemetry config in `lib/go/telemetry` (dev), alert filters in `infra/observability/provisioning/alerting.yml` (infra).

**Effort.** S.

---

# Details — 👨‍💻 Developer (not infra work)

### 11. Error tracking = error logs + trace error spans (zero new infra)

**Problem.** Backend errors today are greppable in logs only — no alert, no grouping, no "did v0.36.0 introduce this 500" answer. Dedicated trackers (Glitchtip, Sentry self-hosted) cost ~512MB–1GB each — real money on an 8GB box with a staging twin coming. Cut.

**Who.**
- **Developer:** `RecordError(ctx, err)` helper in `lib/go/telemetry` + calls at **caught-error sites across the four services** + wire the git tag into `semconv.ServiceVersion`. This is application code.
- **Infra (small tail):** the Grafana panel ("top errors") + the alert rule.

**Fix.** You already run everything this needs; the work is wiring, not new services:
- **Errors land on traces.** `otelhttp` in the harness (`lib/go/httpserver`) already marks 5xx responses as error spans by default. Extend it: a `RecordError(ctx, err)` helper in `lib/go/telemetry` (→ `span.RecordError(err)` + `SetStatus(codes.Error, ...)`), called at **caught-error sites** — business-layer failures that return 4xx or get logged and swallowed. Errors become visible in trace search, not just logs.
- **Error logs carry their trace.** The OTLP logs exporter (#9) records `trace_id`/`span_id` automatically when the log happens inside a span — every error log links to the trace of the request that caused it. Log → trace jump in Grafana.
- **Alerts wake you.** A Grafana rule on error-level log rate → the same ntfy contact point as #2.
- **Grouping & releases without a tracker.** "Top errors" panel via LogsQL `stats by (msg)` over error-level entries = grouping. The tracer's hardcoded `semconv.ServiceVersion("dev")` (in `lib/go/telemetry/otel.go`) gets wired to the git tag, so every error is tagged with the release that produced it — "did v0.36.0 introduce this?" is one label filter.
- **Escalation path:** if this ever feels insufficient, **Sentry cloud free tier** (1 user, ~5k errors/mo, verify at signup) is a DSN swap away — same SDK, zero hosting. The self-hosted instinct is worth dropping; SDK compatibility keeps migration a one-line change either way.

**Where.** `lib/go/telemetry` (+ error sites in `api/*`) — dev; `infra/observability/provisioning/alerting.yml` — infra.

**Effort.** S (after #9). No new containers, ~0MB added.

### 16. E2e package (Playwright) against staging

**Problem.** The staging plan's entire payoff is "e2e tests hit `staging.*` URLs and get prod-like behavior on a live box" — but **no e2e suite exists anywhere in the repos today**. Without it, staging is just a second prod nobody watches.

**Who.** **Developer** — a Playwright suite is test/application work; no infra component (infra's staging work — #12–#18 — is the environment it runs against).

**Fix.** A Playwright suite (`e2e/` in TheTree): signup → login → create event → ticket flow, asserting mailpit (`http://127.0.0.1:8026/api/v1/messages` via SSH tunnel) and TEST_MODE payments. Workflow runs it after `deploy-front-staging` (and nightly after #15). Start with 3–5 happy paths; the point is a heartbeat, not coverage.

**Where.** New `e2e/` in TheTree + workflow.

**Effort.** M. This is the missing piece of the staging plan — schedule it alongside staging itself.

---

# Details — 🧰 Infra DX (rebuild, move & ergonomics)

### 35. Off-box secrets vault (age/sops)

**Problem.** The #1 blocker to ever moving or rebuilding: every secret lives only on the box — deploy's 5 env files, `forgejo/.env`, `caddy/.env`, `observability/.env`, beszel's KEY/TOKEN — plus the CF API token (Forgejo UI) and worker secrets (CF UI). A new box without the old box = can't boot, can't deploy. Nothing off-box exists today.

**Fix.** age-encrypted vault (or sops) with scripted backup/restore:
- `just secrets-backup` → one encrypted file (all envs + a manifest of CF/worker secret *names*) pushed to the same target as #1.
- `just secrets-restore` on a fresh box → all envs back in place.
- The manifest documents **which secrets must be re-created manually in UIs** (CF API token, worker secrets) — the things a file can't restore.
- Key management: one age key, kept off-box (that's the whole point).

**Where.** `infra/secrets/` + backup target.

**Effort.** M. Pairs with #1 (same target) — together they're what makes the box movable.

### 36. `just bootstrap` — whole box from nothing

**Problem.** `setup.sh` covers only caddy + forgejo (confirmed: zero references to observability/beszel/ntfy/mox). **`obs-net` is created by no script** — a fresh box's observability stack won't start. Host deps (docker, user, swap, ufw), crontabs, and beszel's 9-step UI onboarding are all manual. Env template naming is inconsistent (`.env.example` vs `.example.env`), which trips automation.

**Fix.** Extend setup.sh into a full, idempotent `just bootstrap`:
- host deps: docker + compose, user, ssh keys, swap/zram (#5), ufw (#23), fail2ban (#22), unattended-upgrades (#21)
- both networks: caddy-net + **obs-net**
- all 6 stacks up from templates, with unified env naming (`.env.example` everywhere)
- beszel onboarding scripted (replace the 9-step UI dance)
- crontab install (#19)

**Where.** `infra/setup.sh` + `infra/justfile`.

**Effort.** M–L. This is the "install everything from scratch" answer.

### 37. Rebuild runbook doc

**Problem.** The order of a fresh build matters and is written nowhere: images live in Forgejo's registry *on the box* (chicken-and-egg — Forgejo must be up before the deploy stack can pull `git.trieoh.com/trieoh/*`), and DNS must point at the new box first.

**Fix.** A tested, ordered runbook in the infra repo: DNS cutover → host → Forgejo restore → registry → infra stacks → deploy stack → workers → smoke tests. Written from an actual drill (#39), not from theory.

**Where.** `infra/docs/rebuild-runbook.md`.

**Effort.** S.

### 38. Cert provisioning into git

**Problem.** Zero documentation of how the `*.trieoh.com` origin cert (`/etc/caddy/certs`) and mox's `/etc/letsencrypt` are obtained or renewed. A fresh box's TLS is a "figure it out" step.

**Fix.** Document (and script if possible — likely acme.sh/certbot DNS-01 against CF): how each cert is issued, renewed, and where it lands. Feeds the "certs" step of #37's runbook.

**Where.** `infra/README.md` or the runbook.

**Effort.** S.

### 39. Rebuild drill (annual)

**Problem.** A runbook nobody has executed is a fiction.

**Fix.** Once — then yearly — run #37 on a throwaway VPS for a weekend: DNS in, `just bootstrap`, `just secrets-restore`, smoke tests. Fix whatever breaks. The only way to know the box is actually rebuildable.

**Where.** Calendar + a cheap VPS.

**Effort.** M (one weekend), then L as an annual.

---

# Details — 🧰 Infra simplifications (deletion over addition)

### 40. Delete dead Caddyfile config

**Verified:** `api.trieoh.com` has **zero consumers** repo-wide (no frontend, workflow, or doc references it — every app uses its per-service domain). Its `(api_routes)` snippet is used by no other block. Both are dead config. Also delete the "mailpit is dev-only… left as a reminder" comment block (the reminder has served its purpose).

**Effort.** XS.

### 41. Snippet-ify the repeated mail blocks

The Caddyfile has **3× repeated pairs**: an `http://<host>` acme-challenge handler (identical body) + a TLS block with the **same** `/etc/letsencrypt/live/mta-sts.trieoh.com/...` cert path (mail, mta-sts, autoconfig). One `(acme_challenge)` snippet + one `(letsencrypt_tls)` snippet collapse ~30 lines. (Superseded by #42 if certs consolidate — then these blocks disappear entirely.)

**Effort.** S.

### 42. Consolidate cert provisioning to one source

Three cert mechanisms today: Caddy's automatic Let's Encrypt, `/etc/letsencrypt` (mox), and `/etc/caddy/certs` (wildcard). Every mechanism is a renewal path that can silently break (#8 exists because of this). Consolidate to **one**: either everything through Caddy's ACME, or one certbot wildcard shared by all. If mox moves onto the shared source, the hand-rolled acme-challenge blocks and the second cert path die with it (#41 becomes moot).

**Effort.** M. Cross-ref: #8 (cert expiry alert), #38 (cert docs).

### 43. Kill the forgejo-restart cron by fixing its root cause

A **nightly forgejo-restart cron** exists (server crontab, tribal). Nightly restarts are a symptom patch — almost certainly DinD/runner memory growth. Investigate the actual leak (limits #7, runner version, DinD cache growth), fix it, and **delete the cron**. A service that needs nightly restarts is a service that's slowly failing.

**Effort.** M (investigate + fix, then the cron dies).

### 44. Rename forgejo's confusing `internal` network

forgejo's compose declares `internal: driver: bridge, internal: false` — while thetree's compose declares `internal: internal: true`. Same name, opposite meaning (project-scoped so no collision, but anyone reading a compose file gets whiplash). Rename forgejo's to something honest (e.g. `runner-net`).

**Effort.** XS.

### 45. `just register-runner <token>` script

The infra README documents a manual `docker run --rm -it … forgejo-runner register` dance (token from the UI, `--labels` flags, network name). One `just register-runner <token>` recipe wrapping that command makes runner setup a single line and keeps the flags in git.

**Effort.** XS.

### 46. Forgejo composite action — pinned, cached docker CLI ✅ **done 2026-08-24**

**Implemented:** `.forgejo/actions/setup-docker/action.yml` (TheTree) + all 5 workflow call sites swapped; stale runner label in `infra/README.md` fixed.

**Verified:** "Set up Docker" (`curl -fsSL https://get.docker.com | sh`) is copy-pasted in **5 spots across 4 workflows** (ci ×2, frontend-lint-tsc, publish, trivy). Worse, it installs a **full Docker daemon that is never used** — the runner config already automounts the Docker socket into job containers (`docker_host: automount`), so jobs only need the **CLI**. And it's unpinned (`get.docker.com` = latest every run) + re-downloads every run.

**Design:** `.forgejo/actions/setup-docker/action.yml` (local composite action, `uses: ./.forgejo/actions/setup-docker`):

```yaml
name: Setup Docker CLI
inputs:
  version:
    required: false
    default: "28.3.0"   # pinned — bump deliberately, never latest
runs:
  using: composite
  steps:
    - uses: https://code.forgejo.org/actions/cache@v2
      id: cache
      with:
        path: /usr/local/bin/docker
        key: docker-cli-${{ inputs.version }}
    - if: steps.cache.outputs.cache-hit != 'true'
      run: |
        curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-${{ inputs.version }}.tgz -o /tmp/docker.tgz
        tar -xzf /tmp/docker.tgz -C /tmp docker/docker
        install -m 0755 /tmp/docker/docker /usr/local/bin/docker
      shell: bash
    - run: docker version   # smoke — talks to the automounted socket
      shell: bash
```

- **Pinned:** static tarball from `download.docker.com` (no `latest`, no apt) — the version is one input.
- **No download per run:** the runner's internal cache (`forgejo-runner:7080`, `cache.enabled: true`) restores the binary; the tarball downloads once, ever.
- **Less waste:** CLI-only (~25MB) replaces a full daemon install.
- Replace the 5 copies with one `uses:` line each.

**Bonus:** the infra README's register command says `--labels …ubuntu-22.04` but `runner/config.yml` labels are `node:24.15.0-bookworm` — docs are stale; fix while touching this.

**Effort.** S.

---

## Explicitly out of scope (deliberate)

- **Watchtower / Renovate** — pin-in-git is the version ledger (ADR-0007); auto-update fights it. #27 (bump PRs) is the ergonomic middle ground.
- **K8s / Swarm / Traefik** — nothing here needs orchestration; Caddy static config + compose is right-sized for one box.
- **A second server** — the whole point of the staging plan is to fit the twin on the same box.
- **Sentry self-hosted / Glitchtip** — 16GB (Sentry) or ~512MB + own postgres (Glitchtip) doesn't fit an 8GB box with a staging twin; #11 is the zero-infra alternative, and Sentry cloud's free tier is the escalation path.

## Suggested order of operations (the priority guide, regardless of owner)

1. **Eyes on the box (infra):** #2 alerts → ntfy, #3 CF uptime cron, #4 log rotation + disk alert.
2. **Backups (infra):** #1 (with the explicit scope note) — before anything else changes. Include the restore drill. **#35 secrets vault rides along in the same target** — it's what makes the box movable.
3. **Memory (infra):** #5 swap/zram, then #6 profile, then #7 limits everywhere (data-driven).
4. **Surface & secrets (infra):** #8 cert expiry (incl. mail hosts), #22 SSH hardening, #23 host firewall, #24 gitleaks.
5. **Searchable logs + error tracking (shared + dev):** #9 log shipping — **dev's exporter and infra's wiring in the same window** — then #10 latency alerts, then #11 error logs + trace error spans.
6. **Staging launch day (infra):** #12 auth gate + #13 wildcard DNS + #14 mailpit retention + #7 caps on staging, together — and #18 staging telemetry the same day, so staging is never blind.
7. **Staging payoff (dev):** #16 e2e package, then (infra) #15 nightly reset, then #17 drift check.
8. **Steady state (infra):** #19 crontab in git, #20 image GC, #31 pin informd, #28 `--wait`, #29 caddy validate, #30 ops kit + rollback paths.
9. **Whenever (infra):** #32 Tailscale (also resolves #22), #27 bump PRs, #21/#26 optional, #25 DNS-as-code.
10. **Deferred by design (infra):** #33 PgBouncer, #34 PITR.
11. **Infra DX (infra):** #35 secrets vault lands with #1; #36 bootstrap, #37 runbook and #38 cert docs when there's a spare weekend; then #39 drill — the proof it all works.
12. **Infra simplifications (infra):** #40 dead Caddyfile config + #44 network rename + #45 register-runner — XS cleanups whenever; #41/#42 cert consolidation with #38; #43 forgejo-restart root cause when there's a quiet afternoon. **#46 pinned+cached docker CLI** — do first, it touches every CI job you'll edit later.
