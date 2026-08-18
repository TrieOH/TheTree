# TrieOH (TheTree)

Monorepo for **TrieOH** — a SaaS platform of four Go microservices behind a
single Caddy gateway, with React front-ends deployed to Cloudflare Workers.

**This repo is source + CI.** To deploy, go to
[`TrieOH/deploy`](https://git.trieoh.com/TrieOH/deploy).

## Stack

| Layer | Tech | Notes |
|---|---|---|
| Backends | Go 1.26 · chi · sqlc · goose · river | 4 microservices in `api/` |
| Frontends | React 19 · TanStack Start · Tailwind v4 | 4 SPAs in `front/`, deploy to Cloudflare Workers |
| Shared | `lib/go` (authz, crypto, db, telemetry, oauth) · `lib/ts` (orval-generated TS clients) | |
| SDKs | `sdk/go`, `sdk/ts` | IdentityX + Payssage public SDKs |
| CI/CD | Forgejo Actions · Dagger (`.dagger/`) | builds + publishes to `git.trieoh.com/trieoh/<svc>` |
| Dev infra | Docker Compose (`compose.yml`) | postgres, rustfs, mailpit + hot-rebuilt services |

## Repo layout

```
api/<svc>/       # one Go service per dir (cmd/, internal/, db/, api-spec.yml)
front/<svc>/     # one React SPA per dir
lib/go/          # shared Go library
lib/ts/          # orval-generated TS clients (types + TanStack Query hooks)
sdk/go/ sdk/ts/  # public SDKs
.dagger/         # Dagger module (ci, lint, test, compile, publish)
docs/            # CONTEXT.md (domain glossary), adr/, agents/
```

## Prereqs

- Go 1.26 (`go.work` at the root)
- pnpm 10
- Docker + Docker Compose
- `just` — all dev recipes live in the `justfile`
- `golangci-lint` (for `just lint`)
- `gotestsum` (for `just test`) — `go install gotest.tools/gotestsum@latest`

## Quickstart

**Backend:**
```bash
cp .example.env .env                       # root: postgres + rustfs creds
cp api/<svc>/.example.env api/<svc>/.env   # per service (×4: identityx, univents, payssage, informd)
just up               # postgres, rustfs, mailpit + all four services (built locally)
just identityx        # run one service in dev — or: univents, payssage, informd
```

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

Dev ports: postgres `5432` · rustfs `9000/9001` · backends `8080`–`8083`
(pprof `6060`–`6063`) · frontends `3000`, `3001`, `3002`, `3004` · mailpit
`8025`.

## Daily commands

| Task | Command |
|---|---|
| Run dev stack | `just up` / `just down` |
| Run one backend | `just <svc>` (e.g. `just univents`) |
| Run all frontends | `pnpm dev` |
| Run one frontend | `pnpm <svc> dev` (e.g. `pnpm univents dev`) |
| Frontend build | `pnpm -r build` |
| Go tests (all) | `just test` |
| Lint (Go) | `just lint` or `just lint <svc>` |
| Frontend typecheck | `pnpm -r tsc` |
| Frontend lint/format | `pnpm -r lint` / `pnpm -r format` |
| Bump Go deps | `just goup` |

## Codegen (run after changing specs/models)

```bash
just generate-oapi     # oapi-codegen: api-spec.yml → internal/openapi bindings
just generate-orval    # orval: TS client + TanStack Query hooks in lib/ts/<svc>/client
```

Generated code is **not committed** — `internal/openapi/` and
`lib/ts/<svc>/client/` are gitignored and regenerated in CI/builds.

## Releases & deploys

- **Backends:** tag `*/v*` (e.g. `identityx/v0.35.3`) → `publish.yml` runs
  Dagger → pushes `git.trieoh.com/trieoh/<svc>:<tag>` (+ `:latest`) to the
  registry.
- **Frontends:** `main` + `front/**` changes → `deploy-front.yml` → Cloudflare
  Workers (wrangler).
- **TS SDKs:** `publish-ts-sdks.yml` publishes `@trieoh/*` packages.
- **Prod deploy:** images are consumed by `TrieOH/deploy` (server pulls that
  repo, not this one). Version bumps = a commit in the deploy repo.

## Docs

- `docs/` in this repo is for **ADRs, agent plans, and agent context** —
  [`CONTEXT.md`](CONTEXT.md) (domain glossary), [`docs/adr/`](docs/adr/)
  (architecture decisions), [`docs/agents/`](docs/agents/) (agent
  conventions), [`.agents/AGENTS.md`](.agents/AGENTS.md) (agent repo guide).
- Human-facing docs live in [`TrieOH/docs`](https://git.trieoh.com/TrieOH/docs).
