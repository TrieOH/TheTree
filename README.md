# TrieOH (TheTree)

Monorepo for **TrieOH** — a SaaS platform of four Go microservices behind a
single Caddy gateway, with React front-ends deployed to Cloudflare Workers.

**This repo is source + CI + the deploy trigger.** Releases deploy by pushing a
tag here; the `TrieOH/deploy` repo is the pipeline-written ledger (see
[Releases & deploys](#releases--deploys)).

## Stack

| Layer | Tech | Notes |
|---|---|---|
| Backends | Go 1.26 · chi · sqlc · goose · river | 4 microservices in `api/` |
| Frontends | React 19 · TanStack Start · Tailwind v4 | 4 SPAs in `front/`, deploy to Cloudflare Workers |
| Shared | `lib/go` (authz, crypto, db, telemetry, oauth) · `lib/ts` (orval TS clients, `ui-*`, `front-core*`) | framework-free by default, see *Package boundaries* |
| SDKs | `sdk/go`, `sdk/ts` | IdentityX + Payssage public SDKs |
| CI/CD | Forgejo Actions (`deploy.yml`) | tag-gated checks → publish digests to `git.trieoh.com/trieoh/<svc>` → digest-pinned deploys |
| Dev infra | Docker Compose (`compose.yml`) | postgres, rustfs, mailpit + hot-rebuilt services |

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

- Go 1.26 (`go.work` at the root)
- pnpm (pinned via `packageManager` in package.json — corepack picks it up)
- Docker + Docker Compose
- `just` — all dev recipes live in the `justfile`
- `golangci-lint` (for `just lint`)
- `gotestsum` (for `just test`)
- `trivy` (for the pre-push hook's whole-repo scan)

All three are installable with `just setup` (pinned versions). `go install`
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
| Go tests | `just test` (all) · `just test <svc>` · `just <svc> test` |
| Lint (Go) | `just lint` or `just lint <svc>` |
| Frontend typecheck | `pnpm -r tsc` |
| Frontend lint/format | `pnpm -r lint` / `pnpm -r format` |
| Package boundary check | `pnpm check:boundaries` |
| Bump Go deps | `just goup` |

## Codegen (run after changing specs/models)

```bash
just generate-oapi     # oapi-codegen: api-spec.yml → internal/openapi bindings
just generate-orval    # orval: framework-free TS client in lib/ts/<svc>/client
```

Generated code is **partly committed**: `internal/openapi/` is gitignored and
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
  R-chain checks → build once → push `git.trieoh.com/trieoh/<svc>@sha256:…` →
  ledger commit pins the digest in the deploy repo → VPS `compose pull && up -d`
  (only containers whose digest changed are recreated).
- **Frontends:** tag `<app>-ui/v<semver>` → Worker version uploaded, then the
  SAME version promoted to production (build once, never rebuild).
- **TS SDKs:** tag `<sdk>-sdk-ts/v<semver>` → npm publish (version guard:
  tag == package.json version).
- **Rollback:** backends = revert the ledger commit in `TrieOH/deploy` (the
  previous digest redeploys); fronts = `wrangler versions deploy --version-id
  <previous>`.
- **Discipline:** hotfix tags skip most checks — a hotfix must get its next
  normal release promptly (auditable in tag history).
- **Agents/LLMs must never push tags** — tags are prod deploys, humans only.
  See CONTEXT.md.

## Docs

- `docs/` in this repo is for **ADRs, agent plans, and agent context** —
  [`CONTEXT.md`](CONTEXT.md) (domain glossary), [`docs/adr/`](docs/adr/)
  (architecture decisions), [`docs/agents/`](docs/agents/) (agent
  conventions), [`.agents/AGENTS.md`](.agents/AGENTS.md) (agent repo guide).
- Human-facing docs live in [`TrieOH/docs`](https://git.trieoh.com/TrieOH/docs).
