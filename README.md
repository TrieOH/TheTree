# TrieOH (TheTree)

Monorepo for **TrieOH** — a SaaS platform of four Go microservices on
Cloudflare Containers, with front-ends deployed to Cloudflare Workers.

**This repo is source + CI + the deploy trigger.** Releases deploy by pushing a
tag here (see [Releases & deploys](#releases--deploys)).

## Stack

| Layer | Tech | Notes |
|---|---|---|
| Backends | Go 1.27 · chi · sqlc · goose · river | 4 microservices in `api/` |
| Frontends | React 19 · TanStack Start · Tailwind v4 | 4 SPAs in `front/`, deploy to Cloudflare Workers |
| Shared | `lib/go` (authz, crypto, db, telemetry, oauth) · `lib/ts` (orval TS clients, `ui-*`, `front-core*`) | framework-free by default, see *Package boundaries* |
| SDKs | `sdk/go`, `sdk/ts` | IdentityX + Payssage public SDKs |
| CI/CD | GitHub Actions (`deploy.yml`) | tag-gated checks → image pushed once (GHCR) → Cloudflare Containers / Workers |
| Dev infra | Docker Compose (`compose.yml`) | four service Postgres containers, rustfs, mailpit + hot-rebuilt services |

## Repo layout

```
api/<svc>/       # one Go service per dir (cmd/, internal/, db/, api-spec.yml)
front/<svc>/     # one React SPA per dir
lib/go/          # shared Go library
lib/ts/          # shared TS packages: orval clients, ui-react/ui-solid,
                 # front-core (+ -react / -solid bindings)
sdk/go/ sdk/ts/  # public SDKs
docs/            # CONTEXT.md (domain glossary), adr/, agents/
```

### Package boundaries

Every workspace package declares its framework in its **name**:

| Suffix | Meaning | Examples |
|---|---|---|
| *(none)* | framework-free — React/Solid must not appear in `dependencies` | `@trieoh/front-core`, `@trieoh/univents-api`, `@trieoh/api-client` |
| `-react` | React only | `@trieoh/ui-react`, `@trieoh/front-core-react` |
| `-solid` | Solid only | `@trieoh/ui-solid`, `@trieoh/front-core-solid` |

This is what keeps the Solid app from pulling React in (and vice versa): the
React bindings live in `front-core-react`, the Solid ones in `front-core-solid`,
and the shared logic in the neutral `front-core`. `pnpm check:boundaries` (a
plain `node` script, also wired into the pre-commit hook and CI) fails the build
when a package imports across that line.


## Prereqs

- Go 1.27 (`go.work` at the root)
- pnpm (pinned via `packageManager` in package.json — corepack picks it up)
- Docker + Docker Compose
- `just` — all dev recipes live in the `justfile`
- `sqlc` (for `just gen`)
- `golangci-lint` (for `just lint`)
- `gotestsum` (for `just test`)
- `gitleaks` (for the pre-commit secret scan)
- `trivy` (for the pre-push hook's whole-repo scan)

All of them are installable with `just setup`, at the versions pinned in
`tools.versions` — the one file the justfile, the Dockerfiles and CI read, so
a bump happens there and nowhere else. `go install`
puts binaries in `$(go env GOPATH)/bin` and trivy goes to `~/.local/bin` —
add both to `PATH` if a command reports "command not found":

```bash
just setup
export PATH="$(go env GOPATH)/bin:$HOME/.local/bin:$PATH"   # bash / zsh
fish_add_path "$HOME/.local/bin" "$HOME/go/bin"             # fish (already set up in config.fish)
```

## Hooks (enforced locally, not in CI)

Main pushes run **zero CI** — quality is gated at the deploy lane and by local
hooks:

- **`commit-msg`** — validates `<scope>/ACTION: message` against the scope tree
  *discovered from the repo layout*; the declared scope must be the closest
  common ancestor of the staged changes (hard fail, tells you the right scope).
- **`pre-commit`** — gitleaks diff scan of staged changes (secret-only, fast),
  then per-member checks against each member's actual diff (Go: lint +
  `gotestsum -short`; fronts: `check`/`tsc`/`test`; SDKs: typecheck).
- **`pre-push`** — branch naming must match the scope of the *entire* pushed
  commit range. No repo-wide trivy — the deploy lane owns security scanning.

**Never bypass with `--no-verify`** — the deploy lane re-checks at tag time,
but the per-commit secret scan is the only thing scanning your diffs between
releases.

## Quickstart

**Backend:**
```bash
cp .example.env .env                       # root: rustfs creds
cp api/<svc>/.example.env api/<svc>/.env   # per service (×4: identityx, univents, payssage, informd)
just up               # postgres (with four service DBs), rustfs, mailpit + all four services
just run identityx    # run one service in dev — or: univents, payssage, informd
```

**IdentityX** (ADR-0011) — one image, no runtime modes:
```bash
just run identityx         # compose container, jobs in-process, SMTP → mailpit
just migrate identityx     # migrate over the direct (unpooled) endpoint
```
The API always connects through a pooler (`identityx-pooler` locally, Neon's
pooled endpoint when deployed). It deploys to Cloudflare Containers from
`api/identityx/cloudflare` (`pnpm -C api/identityx/cloudflare deploy`), with
jobs on a Cloudflare Queue and Cron Triggers.

**Frontend:**
```bash
pnpm install
cp front/<svc>/.env.example front/<svc>/.env   # fill in real values
pnpm <svc> dev        # one app (e.g. pnpm univents dev) — or pnpm dev for all
```

Env notes (all `*.env` files are gitignored):
- Backend: the example defaults work for local (postgres/rustfs on
  localhost); fill secrets as needed.
- Frontend: the example's defaults already point at the local backends
  (`VITE_API_URL=http://localhost:808x`). Minimum to fill:
  `AUTH_SESSION_PASSWORD` (32+ chars); PostHog/upload keys are optional in dev.
- Auth runs through the app's BFF (`AUTH_TRANSPORT=bff`) — the browser talks
  to TanStack Start server functions, which call the APIs.

Dev database ports: IdentityX `5432`, Univents `5433`, Payssage `5434`, Informd `5435`
(each app gets its own local Postgres container). RustFS `9000/9001` · backends `8080`–`8083`
(pprof `6060`–`6063`) · frontends `3000`, `3001`, `3002`, `3004` · mailpit
`8025`.

Each local Postgres container creates its database and role from its service
env file. Backends never migrate on startup: `just db <svc>` starts the
database and migrates it, `just migrate <svc> [up|status]` migrates on its own
(set `<PREFIX>_DB_HOST`/`_PORT` to point it at another database, e.g. Neon's
direct endpoint).

## Daily commands

Every `just` recipe is `just <command> [target]`; `just` alone lists them.

| Task | Command |
|---|---|
| Run dev stack | `just up` / `just down` |
| Run one backend | `just run <svc>` (e.g. `just run univents`) |
| Start + migrate one database | `just db <svc>` |
| Migrate | `just migrate <svc>` · `just migrate <svc> status` |
| Run all frontends | `pnpm dev` |
| Run one frontend | `pnpm <svc> dev` (e.g. `pnpm univents dev`) |
| Frontend build | `pnpm -r build` |
| Go tests | `just test [<svc>\|lib\|sdk]` (`-short`) · `just test-all [...]` (with integration) |
| Lint (Go) | `just lint [<svc>\|lib\|sdk]` |
| Frontend typecheck | `pnpm -r tsc` |
| Frontend lint/format | `pnpm -r lint` / `pnpm -r format` |
| Package boundary check | `pnpm check:boundaries` |
| Bump Go deps | `just goup` |

## Codegen (run after changing specs/models)

```bash
just gen               # all backends: sqlc + oapi-codegen, then orval
just gen payssage      # one (or several) backends; orval still regenerates all four
```

Generated code is **partly committed**: `internal/sqlc/` and `internal/openapi/` are gitignored and
regenerated in CI/builds, while the orval TS clients under
`lib/ts/<svc>/client/` **are** committed — after regenerating, run
`pnpm -r tsc` (and `node tools/check-package-boundaries.mjs`) before committing.

## Releases & deploys

Three lanes, one rule of thumb: **main is cheap, deploys are the gate,
`-hotfix.*` is the fire escape.**

| Lane | Trigger | Gate |
|---|---|---|
| Dev | push to `main` | zero CI — hooks only. Main is latest-unstable, never prod. |
| Release | tag `<artifact>/v<semver>` | full R-chain: compile → unit → lint → front → trivy → integration |
| Hotfix | tag `<artifact>/v<semver>-hotfix.N` | compile + unit only (≤2 min tag-to-prod) |

**Any dev with push access deploys by pushing a tag** — nothing else to
touch. The pipeline is the gate: a deploy happens only if every check passes,
and the process is fully automated from there (proven live on all four tag
paths: full/express × service/front).

- **Backends:** tag `<svc>/v<semver>` (e.g. `payssage/v0.7.10`) → `deploy.yml`:
  R-chain checks → build once → push `ghcr.io/trieoh/<svc>@sha256:…` → the same
  image pushed to Cloudflare's registry → Worker + Container rollout from
  `api/<svc>/cloudflare` (the service env is the `<SVC>_API_ENV` secret).
- **Frontends:** tag `<app>-ui/v<semver>` → Worker version uploaded, then the
  SAME version promoted to production (build once, never rebuild).
- **TS SDKs:** tag `<sdk>-sdk-ts/v<semver>` → npm publish (version guard:
  tag == package.json version).
- **Rollback:** `wrangler versions deploy --version-id <previous>` — from
  `api/<svc>/cloudflare` for backends, from `front/<app>` for fronts.
- **Discipline:** hotfix tags skip most checks — a hotfix must get its next
  normal release promptly (auditable in tag history).
- **Agents/LLMs must never push tags** — tags are prod deploys, humans only.
  See CONTEXT.md.

## Docs

- `docs/` in this repo is for **ADRs, agent plans, and agent context** —
  [`CONTEXT.md`](CONTEXT.md) (domain glossary), [`docs/adr/`](docs/adr/)
  (architecture decisions), [`docs/agents/`](docs/agents/) (agent
  conventions), [`.agents/AGENTS.md`](.agents/AGENTS.md) (agent repo guide).
- Human-facing docs live in [`TrieOH/docs`](https://github.com/TrieOH/docs).
