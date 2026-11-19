# TheTree

TheTree is TrieOH's monorepo. It holds four products, the code they share, and the pipeline that deploys them.

| Product | What it does |
| --- | --- |
| IdentityX | Login, accounts, organizations, projects and API keys |
| Payssage | Payments through MercadoPago |
| Informd | Forms and their responses |
| Univents | Events: tickets, check-in, badges and certificates |

Each product is a Go API plus a web app. Every app signs in through IdentityX, and Univents takes payments through Payssage. The APIs run on Cloudflare Containers and the apps on Cloudflare Workers.

New here? Read [`docs/onboarding/`](docs/onboarding/), starting with [Start here](docs/onboarding/start-here.mdx).

## What's where

```
api/<api>/    Go APIs: identityx, payssage, informd, univents
front/<app>/  Web apps (React and Solid), plus the dev portal
lib/go/       Shared Go code
lib/ts/       Shared TypeScript: API clients, UI kits, front-core
sdk/          Go and TypeScript SDKs
docs/         ADRs, onboarding docs, plans
e2e/          End-to-end tests
```

A TypeScript package's name says which framework it may use: `-react`, `-solid`, or no suffix for neither. `pnpm check:boundaries` fails when one imports across that line. The [repo tour](docs/onboarding/repo-tour.md) has the rest.

## Running it

You need Go, Node, pnpm, Docker and `just`. `just setup` installs the other tools, at the versions in `tools.versions`.

```sh
cp .example.env .env
for api in identityx univents payssage informd; do
  cp api/$api/.example.env api/$api/.env
done

pnpm install    # packages and git hooks
just gen        # generated code
just up         # databases, the four APIs, file storage, Mailpit
pnpm univents dev
```

The env files need a few values filled in first. [Local setup](docs/onboarding/local-setup.mdx) lists them, along with the ports.

## Everyday commands

`just` on its own lists every recipe.

| To | Run |
| --- | --- |
| Start or stop everything | `just up` / `just down` |
| Run one API (and the ones it calls) | `just run <api>` |
| Run one app | `pnpm <app> dev` |
| Migrate a database | `just migrate <api>` |
| Rebuild generated code | `just gen` |
| Go tests | `just test <api>` |
| Go lint | `just lint <api>` |
| Typecheck the apps | `pnpm -r tsc` |

Generated code (sqlc, oapi-codegen, the TS clients, Workers types) is never committed. Run `just gen` after you change a spec or a query.

## Commits and hooks

Git hooks check your work locally, and pushing to `main` runs no CI, so the hooks are the check.

- Commits look like `feat(payssage): add invoice generation`. The scope must be the closest folder that holds everything you changed. The hook tells you the right one.
- Branches look like `<scope>/<action>/<name>`, for example `payssage/feat/invoices`.
- Each commit is scanned for secrets, and the code you touched is linted and tested.

Don't skip them with `--no-verify`.

## Deploying

Pushing a tag deploys to production. Pushing to `main` deploys nothing.

| Tag | Deploys |
| --- | --- |
| `<api>/v1.2.3` | An API, after the full checks |
| `<app>-ui/v1.2.3` | A web app |
| `pages/<dir>/v1.2.3` | A Pages site |
| `<sdk>-sdk-ts/v1.2.3` | Checks a TS SDK (publishing to npm is by hand) |
| `...-hotfix.N` | A hotfix: compile and unit tests only, so it ships fast |

Every hotfix should be followed soon by a normal release.

To roll back, run `wrangler versions deploy --version-id <previous>` from `api/<api>/cloudflare` or `front/<app>`.

Tags are for people only. Agents never push them.

## More docs

- [`CONTEXT.md`](CONTEXT.md): what the domain words mean
- [`docs/adr/`](docs/adr/): why things are the way they are
- [`docs/onboarding/`](docs/onboarding/): the guides, also shown in the dev portal
