# Frontend DX Wishlist

Status: **open — nothing built yet**. Date: 2026-08-24. Scope: `TrieOH/TheTree` — `front/{identityx,informd,payssage,univents}` · `lib/ts/{api-client,front-core,ui-base,shared-utils}` · pnpm workspace. Owner: **frontend developer**.

Evidence-based: every item cites the actual files. The stack is strong (pnpm catalog, TanStack Start, Tailwind v4, shared `@trieoh/*` packages consumed by all 4 apps) — the gaps are testing, copied components, and config drift.

| # | Item | Effort | Evidence |
|---|------|--------|----------|
| 1 | Vitest + component tests for all 4 apps | M | 4 tests total, all univents, zero component tests; msw provisioned but unused |
| 2 | Move shadcn/ui components into `@trieoh/ui-base` | M | dropdown/dialog/context-menu/breadcrumb ×4, button/card/badge/input ×3, `lib/utils.ts` + `lib/api/fetch.ts` ×4 |
| 3 | Shared t3-env schema in `front-core` | S | 4 `env.ts` = 285 lines; same vars repeated |
| 4 | Shared tsconfig base | S | 4 standalone configs, zero use `extends` |
| 5 | Unify biome configs | XS | 2 variants: identityx `useIgnoreFile:false`, others + `noExplicitAny` |
| 6 | Fix the orval-client gitignore lie | XS | README says "gitignored", git tracks `lib/ts/<svc>/client/` |
| 7 | Pin `latest` catalog entries | XS | `@tanstack/eslint-config: latest`, `@tanstack/devtools-vite: latest` |
| 8 | Remove unused `msw` from workspace trust config | XS | In `onlyBuiltDependencies`/`allowBuilds`, zero usage |
| 9 | Standardize the app skeleton (optional) | S | `server.ts` only in 2/4 apps, `tracing/` only in univents |
| 10 | Component gallery (optional) | M | No storybook; a kitchen-sink route in ui-base instead |

---

### 1. Vitest + component tests for all 4 apps

**Problem.** The frontend test story barely exists: **4 test files total**, all in univents, and they're integration/tracing tests (`server.test.ts`, `tracing/*.test.ts`) — **zero component or unit tests** in any app. identityx/informd/payssage have no vitest config at all. Telling: `msw` (Mock Service Worker — the standard for mocking API calls in component tests) is provisioned in the workspace's `onlyBuiltDependencies` + `allowBuilds`, i.e. someone set up for testing and never wrote the tests.

**Fix.** A shared vitest config in `front-core` (or a `vitest` workspace package): happy-dom/jsdom, one `setup.ts` with the `@trieoh/api-client` fetch stack mocked (msw — finally earning its keep), path aliases from the shared tsconfig (#4). Then per app: smoke tests for the router (each route renders), and component tests for the highest-value UI (auth flows, forms). Start with ~5 per app; the point is a regression net, not coverage.

**Effort.** M.

### 2. Move shadcn/ui components into `@trieoh/ui-base`

**Problem.** The shadcn/ui components were copied **per app** instead of shared: `dropdown-menu`, `dialog`, `context-menu`, `breadcrumb` ×4; `button`, `card`, `badge`, `input`, `label`, `sonner` ×3; plus `shared/lib/utils.ts` and `shared/lib/api/fetch.ts` in all 4. `@trieoh/ui-base` already exists and is consumed by all 4 apps — it's the natural home. A theme or component fix today means editing the same file 3–4 times (or silently letting them drift).

**Fix.** Diff the copies across apps first (shadcn files get customized per-app — merge the intentional deltas), then move the shared set into `ui-base` (re-exporting `cn`, `lib/utils`). Apps import from `@trieoh/ui-base`; per-app overrides stay local. One place to fix, one place to upgrade shadcn.

**Effort.** M. Biggest DRY win on the frontend side.

### 3. Shared t3-env schema in `front-core`

**Problem.** 4 `env.ts` files, 285 lines total, repeating the same shape: API URLs, auth URLs, PostHog, app title. Each app re-declares the common vars with its own zod schema.

**Fix.** `front-core` exports a `commonEnv` schema (or a `createEnv({ app })` builder); apps extend with their per-app vars (univents' `VITE_SUPPORTED_PROVIDERS`, payssage's `VITE_INTENT_TEST_MODE`, etc.). One place for the shared contract (also mirrors the backend: env validation is already centralized per-service pattern).

**Effort.** S.

### 4. Shared tsconfig base

**Problem.** 4 standalone `tsconfig.json` (37–41 lines each), **zero use `extends`** — the compiler options are copy-pasted and will drift (they already differ slightly: payssage/univents 38–41 lines vs identityx 37).

**Fix.** A shared `tsconfig.base.json` (or a `@trieoh/tsconfig` package) with the common options; each app's config extends it and overrides paths/JSX settings. Also gives #1 (vitest) and #3 (front-core) one alias source.

**Effort.** S.

### 5. Unify biome configs

**Problem.** 2 variants: identityx has `useIgnoreFile: false`, the other three have `useIgnoreFile: true` plus a `suspicious.noExplicitAny: warn` rule identityx lacks. Same team, four apps, two lint configs.

**Fix.** One config (the stricter one — `noExplicitAny: warn` is exactly what a growing codebase wants), shared via a base or a single workspace-level `biome.json` with per-app overrides if ever needed.

**Effort.** XS.

### 6. Fix the orval-client gitignore lie

**Problem.** The README states "Generated code is not committed — `internal/openapi/` and `lib/ts/<svc>/client/` are gitignored and regenerated in CI/builds" — but `git ls-files lib/ts/identityx/client` shows the client **is committed** (verified: not in `.gitignore`). So today: spec changes require regenerating AND committing client code by hand, with no CI check that the committed client matches the spec.

**Fix.** Pick one and make the docs + git agree:
- **Commit them properly** (they already are): add a CI step that regenerates via orval and fails on diff — the committed client can never drift from the spec.
- Or **gitignore + regenerate in CI/builds** (the README's original intent): local devs regenerate via `just generate-orval`.

Recommend the first (commit + CI diff check) — the generated client is the contract the apps import; committing it makes changes reviewable.

**Effort.** XS (decision + one CI step).

### 7. Pin `latest` catalog entries

**Problem.** The catalog pins everything except two: `@tanstack/eslint-config: latest` and `@tanstack/devtools-vite: latest`. Floating versions in a repo that already enforces `minimumReleaseAge` + `trustPolicy` — these two are the supply-chain gap.

**Fix.** Pin both to the current versions (same as the other catalog entries).

**Effort.** XS.

### 8. Remove unused `msw` from the workspace trust config

**Problem.** `msw` is in `pnpm-workspace.yaml`'s `onlyBuiltDependencies` + `allowBuilds` but isn't a dependency of any package (verified: no `package.json` references it) and has zero usage.

**Fix.** Remove it from both lists (and from `package.json` if anything still references it) — the build-script trust is permission to run a postinstall that doesn't need to run. It comes back the day #1's tests actually use it.

**Effort.** XS.

### 9. Standardize the app skeleton (optional)

**Problem.** The 4 apps aren't uniform: `server.ts` exists in identityx (5 lines) and univents (30) but not informd/payssage; `tracing/` exists only in univents. Env shapes and route setups differ in ways that are probably incidental, not intentional.

**Fix.** Document the canonical app shape (client-render vs BFF routes, tracing wiring, env contract) — a short `front/README.md` or a template app. Not a refactor; a spec for what a 5th app would copy. (If a 5th app ever lands, this becomes a real scaffold.)

**Effort.** S. Optional.

### 10. Component gallery (optional)

**Problem.** No storybook — component development happens against real pages, which means exercising a component requires navigating to its page.

**Fix.** Skip Storybook (heavy for solo); a **kitchen-sink route** in ui-base (or one app) that renders every `ui-base` component with variants. Cheap, serves as both gallery and the visual smoke test #1 can assert on.

**Effort.** M. Optional — do only if component churn starts to hurt.

---

## Explicitly not frontend-DX scope

- **e2e against staging** — backend wishlist #16 (Playwright); this list is unit/component-level.
- **Visual regression testing** — overkill solo until #1 exists.
- **Per-app theme divergence** — a real issue, but it's design work, not DX; #2 (ui-base) is the mechanical prerequisite.

## Suggested order

1. **XS cleanups:** #6 orval lie (decision + CI diff step), #7 pin `latest`, #8 msw trust removal, #5 biome unify.
2. **#4 tsconfig base** — unblocks #1's vitest config and #3's aliases.
3. **#3 shared env schema** — rides the same "front-core owns the common" wave as #4.
4. **#2 shadcn → ui-base** — biggest DRY win; diff the copies first.
5. **#1 vitest + component tests** — the biggest gap; build on 2–4.
6. **#9 / #10** — optional, whenever.
