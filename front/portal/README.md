# portal

The TrieOH dev portal: docs, services and the state of TheTree for the people who build it. One Worker serves the Solid SPA, its read API over D1 and, as later phases land, the Linear/GitHub/deploy webhooks and the analytics crons.

**Hard rule:** nothing the portal reads may wake a container or a database. No calls to the four APIs, no Postgres, no `/health` polling. Every number comes from outside the services (Linear, GitHub, Cloudflare analytics, Grafana Cloud, Neon's control plane) and is stored in D1 first.

## Layout

Same stack and structure as `front/univents-solid` (Solid 2, TanStack Solid Router, Tailwind 4, `@cloudflare/vite-plugin`).

- `src/server.ts`: `/api/*` (reads D1), `/webhooks/*`, everything else → ASSETS.
- `src/routes/`: Home, Services, Issues, Delivery, Metrics, Docs.
- `src/features/docs/`: docs come from the repo's `docs/` at build time. `docs/adr/*.md` are compiled by MDX (`mdx/remark-adr.ts` lifts the number, title, status and date into `meta`) and rendered through `src/shared/mdx/jsx-runtime.ts`.
- `src/features/issues/`: Linear. `POST /webhooks/linear` checks the HMAC signature and the signed timestamp, drops redeliveries and upserts into D1. An hourly cron (`syncLinear`) pulls what changed since its last run with the read-only key; its first run is the backfill. Issues map to a service through their Linear project. `GET /api/issues` serves the Issues page, the service pages and Home.
- `migrations/`: D1 schema.

## Develop

```sh
pnpm -F portal db:migrate:local   # once, and after adding a migration
pnpm -F portal dev                # http://localhost:3010
pnpm -F portal tsc && pnpm -F portal check && pnpm -F portal test
```

## Secrets

`wrangler secret put` on the `portal` Worker: `LINEAR_API_KEY` (personal key, read-only, team TRI) and `LINEAR_WEBHOOK_SECRET` (the signing secret of the Linear webhook pointing at `https://portal.trieoh.com/webhooks/linear`, events Issues, Cycles and Projects).

## Deploy

Tag `portal-ui/vX.Y.Z` (see `.github/workflows/deploy.yml`): checks, build, `wrangler versions upload`, D1 migrations, promote.
