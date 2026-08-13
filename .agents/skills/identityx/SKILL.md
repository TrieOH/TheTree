---
name: identityx
description: Use the IdentityX service — auth, actors, projects, service API keys, capabilities. Use when working on or calling api/identityx, bootstrapping a fresh environment (setup → projects → svc actors → API keys), creating users/actors or API keys, wiring IDENTITY_X_* envs, or when another service (payssage, univents, informd) needs its identity project configured.
---

# IdentityX

IdentityX is the monorepo's identity/tenant service: actors (humans, services, machines), projects (workspace scopes), organizations, Ed25519-signed JWTs, and service API keys. Every other backend (payssage, univents, informd) bootstraps its own project here and authenticates with a **service API key** that lives on a **svc account actor**.

Run locally: `docker compose up identityx` → `http://localhost:8080` (env in `api/identityx/.env`). The full contract is `api/identityx/api-spec.yml`.

## The actor model (read CONTEXT.md for the domain language)

| Concept | Meaning |
|---|---|
| Actor | human / service / machine identity; everything that can authenticate |
| Project user | actor scoped to a project (`actors.project_id` set, **no role row**) |
| Project member | actor with a `project_members` role (owner/admin/member) |
| Organization | top-level tenant; projects can be org-scoped |
| API key | `{brand_slug}_v1_{env}_{random}`; bound to a subject actor; hash stored, raw returned once |

A **service API key** = an API key whose subject is an actor of `type: service`, `auth_method: api_key`, created in the service's own project. This is what backends put in `IDENTITYX_ACCESS_API_KEY` — not a human's key.

## Bootstrap a fresh environment (do this in order)

All requests are `POST/GET http://localhost:8080` with `Content-Type: application/json`.

```bash
# 1. First account (super admin, platform-level, succeeds once)
curl -s -X POST /auth/setup -d '{"email":"admin@trieoh.com","password":"Admin123#"}'
#  → save data.access_token as $ADMIN_JWT

# 2. One project per service (brand_slug becomes the API-key prefix)
curl -s -X POST /projects -H "Authorization: Bearer $ADMIN_JWT" \
  -d '{"name":"Univents","brand_slug":"univents","domain":"https://univents.com.br"}'   # → $PROJECT_ID

# 3. A svc account actor in that project (auth_method=api_key)
curl -s -X POST /projects/$PROJECT_ID/actors -H "Authorization: Bearer $ADMIN_JWT" \
  -d '{"auth_method":"api_key","type":"service","email":"univents-svc@trieoh.com"}'

# 4. The service API key ON that actor (subject = the svc actor; the caller
#    must be a project admin — the project creator is the owner)
curl -s -X POST /projects/$PROJECT_ID/api_keys -H "Authorization: Bearer $ADMIN_JWT" \
  -d '{"subject_id":"<svc_actor_id>","name":"univents-svc-key","env":"production"}'
#  → data.raw_key is returned ONCE; store it in the service's .env
```

Then set the service's env: `IDENTITY_X_PROJECT_ID=<project>`, `IDENTITY_X_API_KEY=<raw key>`. If the project id changes, **recreate the containers** — `docker compose restart` does not re-read `.env`.

## Everyday operations

- **Register/login a project user** (what frontends do): `POST /auth/register?project_id=…` then `POST /auth/login?project_id=…` → `data.access_token` (JWT) + refresh. The JWT carries `subject.id` (the actor id — used as the purchaser/owner id downstream).
- **Introspect** (what backends do to resolve an API key): `GET /auth/introspect` with `Authorization: Bearer <jwt>` **or** `X-API-Key: <raw key>`. Response is the fun envelope: `data.subject.id`, `data.subject.project_id`, `data.cred.type` (`token`|`api_key`).
- **Members**: `POST /projects/{id}/members` `{actor_email, role}` (owner/admin/member); `POST /organizations/{id}/members` for orgs.
- **Actors**: `GET /projects/{id}/actors`, `GET /projects/{id}/actors/{id}`, `GET /projects/{id}/actors/{email}:by_email` — open to project **users** (svc accounts included, no membership row needed) and project members; platform-level clients pass only as members, never by being platform-level. Not public. `createActor` stays platform-level only.
- **Capabilities**: `GET/POST /projects/{id}/capabilities`; keys can be scoped by capability ids.
- **OAuth providers** (Google, GitHub…): `GET/POST /projects/{id}/oauth-providers`; the browser flow is `/auth/{provider}/connect` → `/auth/{provider}/callback`.

## Gotchas (learned the hard way)

- The API-key route is `/projects/{project_id}/api_keys` (underscore) — `/api-keys` 404s.
- Svc actors are **project users, not project members** — they have no `project_members` row. `projects.GetMember(actor)` misses and the repo maps it to a misleading **"project not found"**. Key creation must rely on the caller's admin `CheckProject`, never a GetMember on the subject (fixed in this repo — keep it that way).
- `createAPIKey` with `subject_id` requires the caller to be a **project admin**; a platform super admin is not automatically a member of a project you created for them — the project creator is added as owner, so the admin JWT that created the project can create keys.
- Introspect **never** returns the bare identity — it's the `{code,data:{subject,cred},…}` envelope. `SetResult(&identity)` on the raw body yields a zeroed identity (owner_id `00000000-…` downstream). Unwrap `data`.
- `POST /auth/setup` 503s payssage/univents boots until it runs: services panic "please setup IDX first".
- Dropping the identityx DB invalidates every API key and project id — recreate keys and re-point envs (they're gitignored).
- Email verification: dev flows may require `/auth/verify-email` or `DEBUG_MODE` skipping it — check the service's env before debugging a "verified" failure.

## Pointers

- `api/identityx/api-spec.yml` — the full contract (authn, actors, projects, orgs, apikeys, capabilities, profiles, email-templates).
- `CONTEXT.md` — Actor / Organization / Project / Project user vs member definitions.
- Sibling skills: `payssage` (its wallet owner is the payssage svc actor), `univents` (its users/buyers are project users here).
