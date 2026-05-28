# TheTree / TrieOH — Project Summary

> Last updated 2026-05-27. Hand-curated; update as the project evolves.

TheTree is the monorepo for **TrieOH**, a SaaS platform built as a suite of domain-specific Go microservices behind a single Caddy API gateway, with React front-ends deployed to Cloudflare Workers. It is a hybrid Go / TypeScript monorepo orchestrated via Docker Compose with self-hosted infrastructure for observability, source control, email, and object storage.

## Repository Layout

| Directory | Language | Purpose |
|-----------|----------|---------|
| `api/` | Go | Four backend microservices (IdentityX, Payssage, Informd, Univents) |
| `front/` | TypeScript/React | Five front-end SPAs targeting Cloudflare Workers |
| `lib/go/` | Go | Shared Go library (authz, crypto, database, telemetry, OAuth, utils) |
| `lib/ts/` | TypeScript | TypeScript type bindings (tygo-generated from Go models) |
| `sdk/go/` | Go | Public SDKs for IdentityX and Payssage |
| `sdk/ts/` | TypeScript | Public TypeScript SDKs for IdentityX and Payssage |
| `docs/` | MDX/TypeScript | Documentation site (Fumadocs, deployed to Cloudflare) |
| `infra/` | Config | Grafana dashboards, VictoriaMetrics scrape configs, Vector pipeline, Docker tooling |
| `bruno/` | JSON | API client collections (Bruno — Postman alternative) |

## Backend Services (Go 1.26)

All services follow a consistent internal structure:
- `cmd/` — entry point
- `internal/app/` — bootstrap, routing, config, dependency injection
- `internal/features/` — domain logic organized per feature
- `internal/database/` — sqlc-generated queries + goose migrations
- `internal/platform/` — infrastructure adapters (SpiceDB, Redis, etc.)
- `internal/shared/` — cross-cutting helpers
- `models/` — domain types with TypeScript generation via tygo

### 1. IdentityX — Auth & Identity Platform
- Actors: human, service, and machine identity types
- Auth methods: password, API key, Google OAuth, GitHub OAuth
- Organizations with role-based membership (member/admin/owner)
- JWT issuance, JWKS endpoint, token introspection, refresh
- Crypto key management, token blacklisting
- Job queue: River; rate-limiting: chi/httprate

### 2. Payssage — Payment Orchestration
- Workspaces with API keys and webhooks
- Payment intents (MercadoPago SDK integration)
- OAuth marketplace connections
- Authorization via SpiceDB
- Job queue: Asynq + Redis; WebSocket support

### 3. Informd — Forms & Notifications
- Multi-tenant namespace isolation
- Form definitions with step-based workflows
- API key management per namespace
- Authorization via SpiceDB
- Job queue: Asynq + Redis

### 4. Univents — Event Management
- Events → Editions → Activities hierarchy
- Checkpoint tracking (attendance)
- Product and ticket management with purchase flow
- Payment integration via Payssage SDK
- Authorization via SpiceDB with granular permissions
- File upload via RustFS (MinIO SDK)
- Job queue: Asynq + Redis; WebSocket support

## Front-End Apps (TypeScript + React 19)

All front-ends use **TanStack Start** (Router + Query), styled with **Tailwind CSS v4**, and deployed to **Cloudflare Workers** via Wrangler. Shared dependencies: `@trieoh/identityx-sdk-ts` for auth, `lucide-react` icons, `motion` animations, `sonner` toasts, `posthog-js` analytics.

| App | Dev Port | Notes |
|-----|----------|-------|
| `front/identityx` | 3000 | BIOME formatter/linter; Radix UI primitives |
| `front/payssage` | 3002 | React Hook Form; Base UI components |
| `front/univents` | 3001 | Leaflet maps; React Day Picker; MercadoPago SDK; Vaul drawers |
| `front/informd` | 3004 | React Hook Form; Base UI; React Compiler (experimental) |
| `front/spicedb` | — | SpiceDB admin UI or local playground |

## Infrastructure (Docker Compose)

Compose files are layered — never run alone:

| File | Role |
|------|------|
| `compose.base.yml` | Foundation: Postgres 18, Caddy, RustFS (S3-compatible storage) |
| `compose.app.yml` | Application services with health-check dependency chains |
| `compose.obs.yml` | Full observability stack: Beszel, VictoriaMetrics/Logs/Traces, Vector, Grafana |
| `compose.prod.yml` | Production: versioned GHCR images, TLS, port exposure |
| `compose.server.yml` | Server infra: Mox (SMTP/IMAP), Forgejo (git forge + CI runners) |
| `compose.dev.yml` | Dev: hot-reload, dev Caddyfile, Mailpit for email testing |

**Startup dependency chain:** Postgres → Caddy → RustFS → IdentityX → (Payssage, Informd, Univents). Univents also depends on Payssage.

## Shared Go Library (`lib/go/`)

- **authz** — SpiceDB client wrapper
- **crypto** — Cryptographic utilities
- **database** — pgx pool, sqlc tx runner, goose migrations, OTel tracing
- **oauth** — OAuth provider abstractions
- **telemetry** — OpenTelemetry setup (traces, metrics)
- **validator** — Input validation
- **utils** — General-purpose utilities
- **xslices** — Generic slice utilities
- **river** — River job queue integration

## Key Technology Choices

- **Go 1.26**, TypeScript 5.7, React 19
- **sqlc** — type-safe SQL code generation from raw SQL
- **tygo** — Go struct → TypeScript type generation
- **SpiceDB / Authzed** — fine-grained authorization (Zanzibar-style)
- **pgx/v5** — PostgreSQL driver
- **Caddy** — reverse proxy with TLS, Cloudflare origin certs
- **OpenTelemetry** — distributed tracing across all services
- **Prometheus** — metrics scraped by VictoriaMetrics
- **River / Asynq** — durable job queues (IdentityX uses River; rest use Asynq)
- **`just`** — command runner

## Development Workflow

- `just dev` — boot observability stack + all four API services + all front-ends
- `just api <svc>` — start specific backend services in Docker
- `just front <svc>` — run specific front-ends via `pnpm dev`
- `just generate` — run tygo across services to regenerate TypeScript types
- `just test <svc>` — run Go tests per service
- `just prod` — deploy production stack with versioned images
- `pnpm -F <pkg> <script>` — run workspace scripts for a specific front-end package
