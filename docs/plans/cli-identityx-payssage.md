# Plan — `idx` & `psx` CLIs (IdentityX + Payssage)

- **Status:** draft for review — no code written yet
- **Date:** 2026-08-25
- **Scope:** `cli/` module (new), `sdk/go/IdentityX`, `sdk/go/Payssage` extensions, `go.work`, `justfile`
- **Owners:** backend developer · platform developer

## Goals

Two command-line faces, covering **every class of user**:

| Who | What they run against | What they need |
|---|---|---|
| Devs running the stack on localhost | local compose (`:8080`/`:8082`) | bootstrap, day-2 ops, testmode |
| Devs operating our prod | prod URLs | same, against prod |
| Clients integrating with our prod | our prod URLs | their project's users, their wallets/intents/webhooks |
| Third parties self-hosting | their own idx/psx URLs | the same tool against *their* deployments |

This is what makes the config model (below) non-negotiable: **environments are
user-defined names for arbitrary deployments**, not TrieOH's own infra. A
client or self-hoster defines their envs and accounts; the defaults
(localhost) are just conveniences.

Per review:
- **One binary per service, both personas inside** — default commands are the
  client surface, `admin …` groups dev ops.
- **Naming:** `idx` / `psx`.
- **Client layer: the public Go SDKs grow** — the CLIs call the SDKs, never
  raw HTTP. Compatible with backend-wishlist #16 (generated Go clients) — when
  that lands, swap SDK internals, CLIs stay put.
- **Config: viper** (`spf13/viper`) — the cobra-native companion; file + env +
  flag precedence with live re-read is its core. We use a scoped subset: one
  YAML file, explicit env bindings, no remote config, no file watching.
- **Accounts & orgs everywhere**: multi-env, multi-account-per-env, live swap,
  default env + default account per env so commands never mention them; org
  commands in both CLIs.

---

## Module layout

**Base CLI framework lives in the shared lib** (`lib/go/cli`, package `cli` in
module `lib`) — every future CLI (informd, univents, …) gets env/account/
login/cfg/output for free. Each actual CLI is a thin binary module under
`cli/` that wires the base root and adds its service commands.

```
lib/go/cli/                  # package cli (module lib) — the BASE every CLI uses
  root.go                    # cobra root factory: base flags, env/account/login
  env.go                     # env add/list/use/show/remove commands
  account.go                 # account add/login/list/use/show/refresh/remove
  cfg.go                     # viper wiring: per-binary config path, env bindings
  output.go                  # JSON printing, secret-once guard, redaction
  ux.go                      # exit codes, arg validation, error footer
  *_test.go                  # unit tests against a temp config (SetConfigFile)

cli/                         # module cli — internal workspace module
                             # (repo convention: relative names like lib, sdk/identityx)
  go.mod                     # module cli
  idx/
    main.go                  cobra root: idx — builds to bin/idx
    auth.go                  setup / login / register / refresh / logout / whoami
    projects.go              projects + members
    actors.go                actors
    apikeys.go               api-keys
    orgs.go                  organizations + org projects/members/actors
    oauth.go                 oauth-providers admin + oauth login
    email.go                 verify-email / resend-verification
  psx/
    main.go                  cobra root: psx — builds to bin/psx
    admin.go                 psx admin … group
    wallets.go               wallets + sandbox + fee + collectors
    intents.go               checkout / get / cancel / refund / list
    webhooks.go              endpoints + events + deliveries
    testmode.go              testmode intent create / refund
    providers.go             connect / revoke
    orgs.go                  organizations + org wallets/intents/collectors/members
  idx/main_test.go, psx/main_test.go   smoke: root + each group renders

go.work                      add ./cli   (lib/go already present)
justfile                     recipes: just cli-build, just cli-test
```

Deps: `spf13/cobra`, `spf13/viper` (+ pflag) in **`lib/go`** (base package) —
stable, widely used, fine alongside the module's existing chi/pgx/zap deps.
The `cli/` binaries themselves only import `lib/cli` + their service SDKs.
SDKs: `resty.dev/v3` replaces `sdkkit` in `sdk/go/IdentityX` (see SDK core).
`env`/`account`/login commands are implemented **once** in `lib/go/cli` and
registered on every root — **one implementation, one config per CLI**
(`idx` writes `idx.yaml`, `psx` writes `psx.yaml`).

---

## Global config model (viper)

**One config per CLI** — each binary owns its own file; nothing is shared:

```
$XDG_CONFIG_HOME/trieoh/idx.yaml       # ~/.config/trieoh/idx.yaml; created 0600
$XDG_CONFIG_HOME/trieoh/psx.yaml       # ~/.config/trieoh/psx.yaml;  created 0600
```

Each file holds only that service's view: its own base URLs and its own
accounts. Same account *names* across the two files = one persona by
convention (e.g. `sophia` in both, with the idx JWT in `idx.yaml` and the
payssage svc key in `psx.yaml`).

`idx.yaml`:

```yaml
default_env: local                       # never mention the env again after setup

environments:
  local:                                 # our compose stack
    idx_url: http://localhost:8080
    default_account: sophia              # default account per env
  prod:
    idx_url: https://identityx.trieoh.com
    default_account: sophia-prod
  acme-dev:                              # a client's self-hosted instance
    idx_url: https://id.acme.dev
    default_account: acme-admin

accounts:
  local:
    sophia:                              # dev persona: platform admin JWT
      token_type: jwt
      access: "<jwt>"
      refresh: "<rt>"
      email: "sophia@trieoh.com"
    alice:                               # a client user on our local stack
      token_type: jwt
      access: "<jwt>"
      refresh: "<rt>"
      project_id: "<acme-proj>"
  prod:
    sophia-prod:
      token_type: jwt
      access: "<jwt>"
      refresh: "<rt>"
    acme-admin:                          # third-party admin on their own deployment
      token_type: jwt
      access: "<jwt>"
      refresh: "<rt>"
```

`psx.yaml` (same shape, service-scoped — envs carry `psx_url` plus an optional
`idx_url` because payssage auth is proxied through identityx, Gaps §2):

```yaml
default_env: local

environments:
  local:
    psx_url: http://localhost:8082
    idx_url: http://localhost:8080        # for api-key cross-check + account refresh
    default_account: sophia
  prod:
    psx_url: https://payssage.trieoh.com
    idx_url: https://identityx.trieoh.com
    default_account: sophia-prod

accounts:
  local:
    sophia:
      token_type: api_key
      api_key: "<payssage svc key>"
      wallet_id: "<platform wallet>"
    alice:                               # same persona as idx.yaml's alice
      token_type: jwt
      access: "<jwt from idx login for payssage project>"
  prod:
    sophia-prod:
      token_type: api_key
      api_key: "<prod svc key>"
      wallet_id: "<prod wallet>"
```

**Model rules:**

- **An account is one credential for this CLI's service**, scoped to one
  environment. The persona spans files by name only — the tool never merges
  the two.
- **Defaults**: `default_env` (top level) + `default_account` (per env) in
  *each* file. After one-time setup, commands never mention env or account.
- **Live swap is per CLI**: `idx env use prod` flips `idx.yaml`'s
  `default_env`; `psx account use sophia --env prod` flips `psx.yaml`'s
  `prod.default_account`. To move a persona, swap it in both files (two short
  commands — acceptable; keeps each config self-contained and independently
  usable by a client who only has one of the two services).
- **Resolution for any command** (viper precedence: **flag > env var > config
  file > built-in default**):
  1. `--env` flag, else `default_env`
  2. `--account` flag, else `environments[env].default_account`
  3. the account block for the credential, env's `*_url` for the base —
     then pushed into the SDK client via `SetCredential`/`SetToken`/
     `SetAPIKey` before the request runs
  4. per-field env-var override remains available for CI/scripts
     (`IDENTITY_X_URL`, `PAYSSAGE_API_KEY`, …) — bound explicitly per key via
     `viper.BindEnv`, not `AutomaticEnv` (we have two env prefixes,
     `IDENTITY_X_*` and `PAYSSAGE_*`, which AutomaticEnv's single prefix
     can't express). Keys are lowercase everywhere.
- **Secrets at rest**: tokens/keys live in the same 0600 file. No keyring —
  dev-tool tradeoff, matches the repo's minimalism. `output` redacts secrets
  in `account show`/`env show` output.
- **File handling**: viper has no atomic save — write via `WriteConfig` +
  `os.Chmod(0600)` on create; a `config init` command scaffolds the file with
  a `local` env and sensible defaults. `cfg` picks the file per binary
  (`idx` → `idx.yaml`, `psx` → `psx.yaml`); tests point viper at a temp file
  via `SetConfigFile`.

### env/account commands (both `idx` and `psx`, same impl, own file)

```
env add <name> --url <base> [--idx-url] [--default-account]
env list                          # active env marked
env use <name>                    # sets default_env in THIS CLI's file
env show [<name>]                 # resolved urls + default account (secrets redacted)
env remove <name>

account add <name> --env <env> --api-key <key> [--project-id] [--wallet-id] [--set-default]
account add <name> --env <env> --token <jwt>   [--project-id] [--wallet-id] [--set-default]
account login <name> --env <env> [--email --password | --provider google|github]   # idx only
account list [--env]              # accounts in env; active marked; cred type shown
account use <name> [--env]        # sets env's default_account
account remove <name> --env
account show [--env] [--account]  # resolved cred for THIS binary's service (redacted)
account refresh [--env] [--account]
```

Notes:
- `idx env add prod --url https://identityx.trieoh.com`; `psx env add prod
  --url https://payssage.trieoh.com --idx-url https://identityx.trieoh.com`
  (psx needs the env's `idx_url` for auth proxying, Gaps §2).
- Env names are per file — they only line up because both sides use the same
  deployment names; the tool doesn't enforce it.

**Login stories:**

- `idx account add --api-key <key>` → validates immediately via
  `GET /auth/introspect` (accepts `X-API-Key`) → stores identity (project,
  actor type, capabilities) with the key. **Works today.**
- `idx account login --email --password` → `POST /auth/login?project_id=` →
  token pair stored.
- `idx account login --provider google|github` → OAuth flow (see Gaps §1 —
  needs API work for a clean loopback flow; paste-back works today).
- `psx account add --api-key <key>` → **no public introspect endpoint on
  payssage** (see Gaps §2); validate lazily on first command, or introspect
  against the env's `idx_url` (the payssage svc key's subject lives in the
  identityx payssage project — same key works).
- `psx account add --token <jwt>` → paste the JWT from `idx login` for the
  payssage project; optionally `--wallet-id`.
- `account refresh` → rotates the stored refresh token via identityx refresh
  (works for psx accounts too, since the JWT and refresh token are
  identityx's — requires the env to have `idx_url`).

**Secrets-once + accounts:** `idx admin api-keys create` and
`psx admin webhooks endpoints add` get a `--save-account <name>` flag that
stores the returned raw key/secret straight into the account block — the only
moment it exists, and the natural place for it.

---

## SDK extension plan (the real work)

### Client core: resty + mutable credentials (both SDKs)

**`sdkkit` is dropped for `resty.dev/v3`** (already pinned in the services at
`v3.0.0-rc.3`; oauth_providers injects a `*resty.Client` today). The IdentityX
SDK's embedded `*sdkkit.Client` is replaced by the same resty-based core the
Payssage SDK gets — one envelope-unwrap helper, one credential model. This
aligns with backend-wishlist #16/#17 (generated Go clients target resty; a
shared `lib/go/httpclient` factory can be adopted by both SDKs when it lands).

**Credentials are mutable, set/get — no `NewWithToken`, no per-call
overrides.** The CLI resolves the active account per command and pushes the
credential into a client; the same client can be reused across accounts:

```go
type Credential struct {
    Type  CredType // jwt | api_key
    Value string
}

func (c *Client) SetCredential(cred Credential) *Client
func (c *Client) Credential() Credential
func (c *Client) SetToken(tok string) *Client      // jwt convenience
func (c *Client) SetAPIKey(key string) *Client     // api_key convenience
func (c *Client) SetProjectID(id uuid.UUID) *Client // per-account project scope (idx)
func (c *Client) ProjectID() uuid.UUID
```

- Requests read the credential **at call time**: `api_key` →
  `X-API-Key: <key>`, `jwt` → `Authorization: Bearer <tok>` — set per request
  in the shared `do()` (`r.SetHeader(...)` on the resty request, no hooks
  needed). `Introspect` accepts either, so it just sends whatever is set.
- `Config` no longer hard-requires `ProjectID`/`APIKey` (today `NewClient`
  errors) — the CLI sets them per account. `BaseURL` stays in `Config` (and
  resty's `SetBaseURL` is available if a client must hop envs).
- Not thread-safe for concurrent credential swaps — the CLI is one command at
  a time, so set-before-use is fine; document it.
- `go.mod` changes: `sdk/go/IdentityX` drops `github.com/MintzyG/sdkkit`,
  both SDKs add `resty.dev/v3`.

Add the missing methods below; keep signatures aligned with `api-spec.yml`.
New code gets `*_test.go` in the existing `httptest` table style
(`sdk/go/Payssage/wallets_test.go` — resty plays fine with `httptest`).

### `sdk/go/IdentityX`

| Method (new unless noted) | Endpoint |
|---|---|
| `Setup(ctx, email, password) (*UserTokens, error)` | `POST /auth/setup` |
| `SetupStatus(ctx) (bool, error)` | `GET /auth/setup` (204 pending / 409 done) |
| `Register(ctx, projectID, email, password) (*UserTokens, error)` | `POST /auth/register?project_id=` |
| `Login(ctx, projectID, email, password) (*UserTokens, error)` | `POST /auth/login?project_id=` |
| `Refresh(ctx, refreshToken) (*UserTokens, error)` | `POST /auth/refresh` |
| `Logout(ctx) error` | `POST /auth/logout` |
| `Introspect(ctx, tokenOrKey) (*Identity, error)` | `GET /auth/introspect` (bearer *or* X-API-Key) |
| `VerifyEmail(ctx, token) error` | `POST /auth/verify-email` |
| `ResendVerification(ctx, projectID, email) error` | `POST /auth/resend-verification` |
| `OAuthConnect(ctx, provider, projectID) (string, error)` | `GET /auth/{provider}/connect` → `data.url` |
| `OAuthCallback(ctx, provider, code, state) (*UserTokens, error)` | `GET /auth/{provider}/callback` |
| `CreateProject(ctx, req) (*Project, error)` | `POST /projects` |
| `ListProjects(ctx) ([]Project, error)` | `GET /projects` |
| `ResolveProjectByDomain(ctx, domainOrURL) (*Project, error)` | `QUERY /projects/by_domain` (RFC 10008) — **new endpoint, hand-wired route (Gaps §7)** |
| `CreateActor(ctx, projectID, req) (*Actor, error)` | `POST /projects/{id}/actors` |
| `ListActors(ctx, projectID) ([]Actor, error)` | `GET /projects/{id}/actors` |
| `AddMember(ctx, projectID, email, role) (*Member, error)` | `POST /projects/{id}/members` |
| `ListMembers(ctx, projectID) ([]Member, error)` | `GET /projects/{id}/members` |
| `CreateAPIKey(ctx, projectID, req) (*APIKey, error)` — **raw key in response** | `POST /projects/{id}/api_keys` |
| `ListCapabilities(ctx, projectID)` | `GET /projects/{id}/capabilities` |
| **Orgs:** | |
| `ListOrganizations(ctx) ([]Org, error)` | `GET /organizations` |
| `CreateOrganization(ctx, name, slug) (*Org, error)` | `POST /organizations` |
| `ListOrganizationMembers(ctx, orgID)` | `GET /organizations/{id}/members` |
| `AddOrganizationMember(ctx, orgID, email, role)` | `POST /organizations/{id}/members` |
| `RemoveOrganizationMember(ctx, orgID, email)` | `DELETE /organizations/{id}/members` |
| `ListOrganizationProjects(ctx, orgID)` | `GET /organizations/{id}/projects` |
| `CreateOrganizationProject(ctx, orgID, req)` | `POST /organizations/{id}/projects` |
| `ListOrgProjectActors(ctx, orgID, projectID)` | `GET /organizations/{org}/projects/{proj}/actors` |
| `CreateOrgProjectActor(ctx, orgID, projectID, req)` | `POST …/actors` (platform-only; same as `createActor`, org-scoped route — CLI uses `admin actors create`; this SDK method optional) |
| `ListOrgProjectMembers(ctx, orgID, projectID)` | `GET /organizations/{org}/projects/{proj}/members` |
| `AddOrgProjectMember(ctx, orgID, projectID, email, role)` | `POST …/members` |
| **OAuth provider admin (needed to register CLI loopback URLs):** | |
| `ListProjectOAuthProviders(ctx, projectID)` | `GET /projects/{id}/oauth-providers` |
| `CreateProjectOAuthProvider(ctx, projectID, req)` | `POST /projects/{id}/oauth-providers` |
| `UpdateOAuthProvider(ctx, id, req)` | `PATCH /oauth-providers/{id}` |
| `DeleteOAuthProvider(ctx, id)` | `DELETE /oauth-providers/{id}` |
| `SetOAuthProviderEnabled(ctx, id, on)` | `POST /oauth-providers/{id}/enable\|disable` |

Existing `ActorService.GetByEmail/GetByID`, `TokenService` stay as-is
(rewritten on the resty core, behavior unchanged).
**Envelope caveat:** the rewrite is the chance to fix the known wart —
`GetByEmail` currently decodes the whole body as the actor instead of
unwrapping `data`. All methods go through one `do()` that unwraps the fun
envelope. Key-creation stays admin-checked on the caller, never a member
lookup on the subject (skill gotcha).

### `sdk/go/Payssage`

| Method (new unless noted) | Endpoint |
|---|---|
| `SetWalletSandbox(ctx, walletID, on bool) error` | `PATCH /wallets/{id}/sandbox` |
| `CreateWebhookEndpoint(ctx, walletID, name, url) (*WebhookEndpoint, error)` — **secret in response** | `POST /wallets/{id}/webhooks/endpoints` |
| `ListWebhookEndpoints(ctx, walletID)` | `GET /wallets/{id}/webhooks/endpoints` |
| `DeleteWebhookEndpoint(ctx, endpointID)` | `DELETE /webhooks/endpoints/{id}` |
| `ListWebhookEvents(ctx, walletID)` | `GET /wallets/{id}/webhooks/events` |
| `ListDeliveries(ctx, endpointID)` | `GET /webhooks/endpoints/{id}/deliveries` |
| `HardCreateIntent(ctx, req) (*Intent, error)` — testmode | `POST /testmode/intents/create` |
| `HardRefundIntent(ctx, intentID)` — testmode | `POST /testmode/intents/{id}/refund` |
| `ListIntents(ctx, walletID)` | `GET /wallets/{id}/intents` |
| `ListWallets(ctx)` | `GET /wallets` (bearer only) |
| `RevokeProvider(ctx, provider)` | `POST /providers/{provider}/revoke` |
| **Orgs:** | |
| `ListOrganizations(ctx) ([]Org, error)` | `GET /organizations` |
| `CreateOrganization(ctx, name, slug) (*Org, error)` | `POST /organizations` |
| `ListOrganizationMembers(ctx, orgID)` | `GET /organizations/{id}/members` |
| `AddOrganizationMember(ctx, orgID, email, role)` | `POST /organizations/{id}/members` |
| `RemoveOrganizationMember(ctx, orgID, email)` | `DELETE /organizations/{id}/member/{email}` |
| `GetOrganizationMember(ctx, orgID, byID *uuid.UUID, byEmail *string)` | `GET …/member/{id}` / `…/{email}:by_email` |
| `ListOrganizationWallets(ctx, orgID)` | `GET /organizations/{id}/wallets` |
| `ListOrganizationIntents(ctx, orgID)` | `GET /organizations/{id}/intents` |
| `ListOrganizationCollectors(ctx, orgID)` | `GET /organizations/{id}/collectors` |

Existing `wallets.go`/`intents.go`/`oauth.go` methods stay as-is (moved to
the shared resty core + credential set/get). `Introspect` is IdentityX-only;
payssage has none (Gaps §2) — psx credentials are validated lazily or via the
env's `idx_url`.

---

## Command trees

### `idx`

```
idx
├── setup --email --password            # first account; prints token pair
│   └── status                          # 204 pending / 409 done
├── login --project-id --email --password   # stores into active/--account
├── register --project-id --email --password
├── verify-email --token
├── resend-verification --project-id --email
├── refresh                              # rotate stored refresh token
├── logout                               # clear stored token
├── whoami [--token]                     # introspect → subject.id, type, project
├── env …                                # shared impl, per-CLI file (idx.yaml)
├── account …                            # shared impl, per-CLI file (idx.yaml)
├── projects create --name --brand-slug [--domain] [--org-id]
├── projects list
├── projects resolve --domain <domain-or-url>      # QUERY /projects/by_domain (Gaps §7)
├── actors list --project-id
├── actors get --project-id --by-email|--id
├── orgs create --name --slug
├── orgs list
├── orgs members list --org-id
├── orgs members add --org-id --email --role owner|admin|member
├── orgs members remove --org-id --email
├── orgs projects list --org-id
├── orgs projects create --org-id --name --brand-slug [--domain]
├── orgs projects actors list --org-id --project-id
├── orgs projects members list --org-id --project-id
├── orgs projects members add --org-id --project-id --email --role
└── admin
    ├── login --email --password         # platform JWT → account (no project scope)
    ├── bootstrap --name --brand-slug --domain --env   # THE dev flow (below)
    ├── projects create|list|get
    ├── actors create --project-id --type service|human|machine \
    │        --auth-method api_key|password [--email]
    ├── api-keys create --project-id --subject-id --name --env \
    │        [--capabilities …] [--expires-at] [--save-account <name>]   # RAW KEY once
    ├── members add|list --project-id
    ├── capabilities list --project-id
    ├── oauth-providers list --project-id
    ├── oauth-providers add --project-id --provider --client-id --client-secret --callback-url
    ├── oauth-providers update --id [--client-id --client-secret --callback-url]
    ├── oauth-providers enable|disable --id
    └── oauth-providers delete --id
```

**`idx admin bootstrap`** — one-shot fresh-env setup (skill flow verbatim):
`GET /auth/setup` → `POST /auth/setup` if pending → `POST /projects` → svc
actor → api key (subject = svc actor) → raw key printed once + env block.
Idempotent: if setup done and brand-slug project exists, only emit env block —
safe to re-run after a DB drop.

### `psx`

```
psx
├── wallet create --name [--organization-id]   # svc key or user JWT
├── wallet get [--wallet-id]
├── wallet list                                 # bearer only (spec quirk)
├── wallet fee --bps [--wallet-id]
├── wallet sandbox on|off [--wallet-id]
├── wallet sellers list [--wallet-id]
├── checkout --wallet-id --amount-cents --method pix|card \
│        [--seller-id --currency]               # prints QR / provider_data
├── intents get --id
├── intents cancel --id
├── intents refund --id
├── intents list [--wallet-id]
├── webhooks endpoints list [--wallet-id]
├── webhooks events list [--wallet-id]
├── webhooks deliveries list --endpoint-id
├── providers connect --provider mercadopago    # prints browser URL
├── providers revoke --provider
├── env …                                      # shared (same impl as idx)
├── account …                                  # shared impl, per-CLI file (psx.yaml)
├── orgs create --name --slug
├── orgs list
├── orgs members list --org-id
├── orgs members add --org-id --email --role
├── orgs members remove --org-id --email
├── orgs members get --org-id --by-id|--by-email
├── orgs wallets list --org-id
├── orgs intents list --org-id
├── orgs collectors list --org-id
└── admin
    ├── bootstrap --name --fee-bps [--webhook-url]   # THE dev flow (below)
    ├── wallet create|get|fee|sandbox                # svc-key variant
    ├── webhooks endpoints add --wallet-id --name --url \
    │        [--save-account <name>]                 # SECRET once
    ├── webhooks endpoints delete --id
    ├── testmode intent create --wallet-id --status succeeded|failed|… \
    │        --amount-cents [--provider --seller-id]  # no provider contacted
    └── testmode intent refund --id
```

**`psx admin bootstrap`** — wallet with the **payssage svc key** (owner = svc
actor, per CheckWalletAccess) → sandbox → fee → optional webhook endpoint →
secret printed once + env block.

---

## Gaps to fix / alter / add (what the API needs for the CLI to work cleanly)

1. **OAuth login needs a loopback story (identityx).** Today the provider
   redirects to the project's stored `callback_url`, which points at
   identityx's own `/auth/{provider}/callback` — the *browser* receives the
   token JSON and someone must paste it back. Works with zero changes
   (paste-back), but a real CLI flow wants the browser to hit
   `http://127.0.0.1:<port>/callback` on the CLI's local listener. To do that:
   - **Add**: allow **multiple callback URLs per provider row** (or a second
     provider row per provider — the repo currently assumes one row per
     provider per project: `GetByProjectAndProvider`, and the exchange uses
     the stored `callback_url` as `redirect_uri`). Google allows multiple
     redirect URIs; **GitHub OAuth apps allow exactly one** — a loopback flow
     can't share GitHub's slot. Minimum viable: support multiple URLs per row,
     document the GitHub caveat (paste-back for GitHub).
   - **Add**: connect-time awareness of which callback URL is in flight so the
     `oauth2.Exchange` uses the same `redirect_uri` the consent URL carried.
   - The CLI side then: `admin oauth-providers add --callback-url http://127.0.0.1:39451/callback`,
     start listener → `idx account login --provider google` → open browser →
     listener catches `?code&state` → `OAuthCallback` → store tokens.
   - State already fits: one-time, 10-min TTL, project-scoped. No change needed.

2. **Payssage auth goes through identityx — by design, not a gap.** psx has
   no login/introspect/refresh of its own; its user JWTs come from identityx
   (the payssage project), and api keys live on the payssage svc actor there.
   The CLI just needs to handle it (no API change):
   - `psx account add --token <jwt>` — paste the JWT from `idx login` for the
     payssage project.
   - `psx account add --api-key <key>` — validate lazily on first command, or
     cross-check against the env's `idx_url` (same key introspects there).
   - `account refresh` on psx accounts calls identityx refresh via the env's
     `idx_url` — psx accounts need the env to have `idx_url`.

3. **Let the server answer auth failures.** No client-side cred gating, no
   `requiredCred` table, no pre-flight checks — the CLI sends the active
   account's credential and the server's 401/403 **is** the answer. The CLI
   renders the envelope error as-is, hints at the fix (`try `idx admin
   login`` or `account use <name>`) in the error footer, and exits 2.

4. **`psx wallet list` is bearer-only and returns 201** (spec quirk). The svc
   key can't enumerate wallets — `admin` works by id (`wallet get`), and the
   account block's `wallet_id` covers the platform wallet. No API change;
   document it.

5. **Verification emails in dev.** `idx register` may hit a verified-email
   gate (or DEBUG_MODE skips it). The CLI should surface the exact error and
   point at `idx verify-email --token` / `idx resend-verification`.
   Optionally: on `register` with `--wait-verify`, poll until verified.

6. **Raw keys are returned once.** `api-keys create` and
   `webhooks endpoints add` must print the secret inside a
   `!!! store this now — it will not be shown again !!!` block. The
   `--save-account` flag (above) is the ergonomic escape hatch. Never log the
   raw key to stderr/debug output.

7. **`idx account login --provider` needs a project scope — and projects
   should resolve by domain/url.** Without a project, `connect` uses platform
   env creds (`oauth.EnvCredentials`); with one, it uses the project's
   provider row. Since `project_id` is a chore to remember, resolve it:
   - **Add (identityx API): a by-domain project lookup using the HTTP QUERY
     method (RFC 10008).** `domain` is already a `UNIQUE` column
     (`uniq_verified_domain`); there is no endpoint or repo method today.
     Add `QUERY /projects/by_domain` with `{"domain": …}` in the body,
     backed by a new `GetByDomain` repo method + sqlc query on the unique
     column.
   - **How it fits the stack (verified empirically + upstream status):**
     - chi v5.3.1 supports QUERY natively — `r.Query(pattern, handler)` and
       `Method("QUERY", …)` (const `methodQuery = "QUERY"` in tree.go).
     - oapi-codegen v2.8.0 (latest) **silently drops `query:` operations** —
       tested: the spec parses, but no interface method, route, or client
       code is generated. So the op goes in the spec as documentation and is
       **excluded via `exclude-operation-ids`** (same mechanism as
       `getHealth`/`getOpenAPISpec`) and hand-registered in
       `httpserver.NewRouter`.
     - **Upstream is actively landing this** (checked 2026-08-25):
       kin-openapi merged OAS 3.2 `query` support (getkin/kin-openapi#1240,
       2026-08-05); the same contributor is doing oapi-codegen's side
       (oapi-codegen#2436, in progress, no release yet — v2.8.0 predates the
       kin-openapi merge). **Migration note**: when a QUERY-capable
       oapi-codegen ships, bump the pinned version (`justfile`
       `oapi-codegen-version`), drop the `exclude-operation-ids` entry and
       the hand-wired route, and regenerate. Chi may need a bump + global
       init per the upstream thread — verify at that time.
     - Go stdlib: `http.NewRequest("QUERY", …)` works with the literal —
       `http.MethodQuery` is still pending upstream (golang/go#80058,
       golang/go#80134), so use the string, not a constant.
   - **CLI**: `--project-id` flags (login/register/account add/api-keys/…) also
     accept a domain or URL and resolve it via that endpoint; convenience
     `idx projects resolve --domain univents.com.br` prints the project.
   - OAuth: `idx account login --provider` requires `--project-id` (or
     `--domain`) unless `--platform` is passed, and warns when the provider
     isn't configured (pre-check with `admin oauth-providers list`).

---

## Output

- Default: pretty-printed JSON of the envelope's `data` (SDKs unwrap it).
- Secrets printed exactly once (see Gaps §6).
- Exit codes: 0 ok, 1 CLI/usage, 2 API error (passes through
  `SDKError`/`APIError`). Errors to stderr.
- Human tables (wallets, intents, keys) are a follow-up; JSON-first keeps MVP
  deps at cobra+viper.

## Testing

- SDK methods: `httptest` table tests mirroring `sdk/go/Payssage/wallets_test.go`
  (envelope fixtures, status branches).
- CLI: `cli/idx/main_test.go` + `cli/psx/main_test.go` — each command against a
  stubbed client interface (thin interface in `lib/go/cli` over the SDKs),
  asserting rendered output/exit codes. `account`/`env` logic lives in
  `lib/go/cli` and is unit-tested there against a temp config file (viper with
  `SetConfigFile`).
- Bootstrap + OAuth paste-back: integration tests behind `IDX_E2E=1` /
  `PSX_E2E=1` against the compose stack.
- `just test cli` covers `./cli/...` + `./lib/go/cli/...`; `just lint`
  likewise (both wired into the existing recipes).

## Build & publish

**No `go install`, no module publishing — we compile and hand out binaries.**
The module path stays `module cli` (repo convention); nobody fetches it.

- Add `./cli` to `go.work` (for workspace builds/tests only).
- `just cli-build` → `(cd cli && go build -o bin/idx ./idx && go build -o
  bin/psx ./psx)` — artifacts in `cli/bin/`.
- **Installing the bins** — the everyday "get the tools" recipe:
  `just cli-install` → builds then copies `bin/idx` + `bin/psx` into
  `~/.local/bin` — the normal place for CLIs on this machine, alongside the
  other tools (`just setup` already puts trivy there; `~/.local/bin` is on
  PATH per the README). Not GOPATH/bin — these aren't go-installed tools,
  they're shipped binaries. Handing binaries to someone else = copy those
  two files; they're static-ish Go builds, no runtime deps beyond libc.
- If distribution is ever wanted (e.g. CI artifacts / release attachments),
  the justfile recipe is the single source — build once, ship the bins. No
  `go install` story at all.
- `just generate` untouched — SDKs are hand-written.

## Phasing

1. **SDK groundwork** — identityx: auth (setup/login/register/refresh/introspect/
   verify/resend) + projects/actors/members/api-keys + orgs + oauth-provider
   admin; payssage: sandbox/webhooks/testmode/intents list + orgs.
   resty core + credential set/get on both. Tests green.
2. **Base CLI + skeleton** — `lib/go/cli` (viper cfg, `env`/`account`
   commands, defaults, live swap, output), cobra roots in `cli/idx` +
   `cli/psx`.
3. **Client commands** — auth/persona commands, wallets/intents/webhooks,
   orgs on both.
4. **Admin groups** — bootstraps, testmode, oauth-provider admin (incl.
   loopback registration).
5. **justfile wiring + E2E.**

## Open questions for review

1. **OAuth scope** — is loopback-callback support (Gaps §1) in scope for this
   pass, or paste-back MVP? GitHub's single-redirect-URI constraint makes this
   a product call, not just plumbing.

## Non-goals (this pass)

- No TS CLIs, no informd/univents CLIs (same pattern applies when wanted).
- No SDK generation refactor (wishlist #16) — SDK-extension choice stays
  compatible with it.
- No keyring, no interactive TUI, no watch/streaming, no
  forgot/reset-password commands (out of MVP scope).
