# Repo tour

- **Order:** 3
- **Summary:** The folders in TheTree, and what is inside an API.

## Top level

| Folder | What's in it |
| --- | --- |
| `api/` | The four Go APIs: `identityx`, `payssage`, `informd`, `univents` |
| `front/` | The web apps, one folder each |
| `lib/go/` | Go code the APIs share: the HTTP server, database, jobs, telemetry |
| `lib/ts/` | TypeScript packages the apps share: API clients, UI components |
| `sdk/` | Public SDKs for IdentityX and Payssage, in Go and TypeScript |
| `docs/` | Decisions (`adr/`) and these guides (`onboarding/`) |
| `e2e/` | End-to-end tests |
| `tools/` | Scripts the build and the git hooks use |

The `justfile` holds the API commands. Run `just` to list them.

`tools.versions` pins the version of every tool except pnpm, which `package.json` pins.

## Inside an API

Every API has the same layout. Using `api/identityx`:

| Path | What it is |
| --- | --- |
| `api-spec.yml` | The OpenAPI spec. Routes, request and response types start here. |
| `problems.yml` | The errors this API can return |
| `cmd/identityx/` | `main`, the entry point |
| `internal/app/` | Startup: config, wiring, and the error mapper |
| `internal/features/<name>/` | One folder per feature, such as `projects` or `api_keys` |
| `db/migrations/`, `db/queries/` | SQL. sqlc turns the queries into Go. |
| `models/` | Shared Go types |
| `errs/` | Errors a client can cause |
| `cloudflare/` | The Worker that runs this API in production |
| `Dockerfile` | The image the Worker runs |

A feature folder holds `handlers.go`, its database code in `repo.go`, and a `service/` folder with one file per operation. A file named after the feature, such as `projects.go`, connects them.

## Inside `front/`

| Folder | What it is | Built with |
| --- | --- | --- |
| `identityx`, `payssage`, `informd` | The product apps | React, TanStack Start |
| `univents-solid` | The Univents app | Solid |
| `univents` | The old Univents app, replaced by `univents-solid` | React, TanStack Start |
| `portal` | This portal | Solid |
| `trieoh` | The TrieOH company site | Solid |
| `trieoh-style-guide` | The design system, at `styles.trieoh.com` | Solid |

The pnpm package for `front/univents-solid` is named `univents`, so `pnpm univents dev` runs the Solid app.

## Shared TypeScript packages

A package's name says which framework it may use:

- a `-react` suffix: React only
- a `-solid` suffix: Solid only
- no suffix: neither

`pnpm check:boundaries` fails when a package breaks this rule.

The ones you'll meet first:

- `lib/ts/<api>/`: the generated client for that API
- `lib/ts/api-client/`: the fetch code every client goes through
- `lib/ts/ui-react/`, `lib/ts/ui-solid/`: UI components

## Generated code

Git ignores code that is built from the specs and the SQL:

- `internal/sqlc/` and `internal/openapi/` in each API
- `api-spec.gen.yml` in each API
- the clients under `lib/ts/<api>/client/`
- every `worker-configuration.d.ts`

`just gen` rebuilds all of it. If an import is missing after a pull, run `just gen`.
