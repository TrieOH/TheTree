# TrieOH — Domain & Architecture Context

Single-context glossary for TheTree. Architecture decisions live in `docs/adr/`.
Terms here give names to good seams; keep this current when concepts are added
or sharpened.

## Domain

- **Actor** — a human, service, or machine identity in IdentityX; everything that can authenticate.
- **Organization** — top-level tenant in IdentityX; members hold roles (member/admin/owner).
- **Project** — an IdentityX workspace under an organization; scope for API keys and profiles.
- **Project user vs Project member** — a project user is an actor scoped to the project (`actors.project_id` set, no role row); a project member holds a `project_members` role (member/admin/owner). Profiles and profile schemas are public reads (anonymous GETs are served); project users may update only their own profile; members may read any profile in the project, and admin/owner may also update any profile in it.
- **Profile / Profile schema** — a user's public identity in IdentityX (name, pfp, socials minimum); shaped by a versioned JSON schema per platform/project; one instance per user carrying its schema version, auto-migrated when a new version validates, else kept and flagged for admin.
- **Event** — a Univents gathering; owned by members with roles (owner/admin/staff).
- **Edition** — one run of an Event, with its own dates and registration window.
- **Program / Program occurrence** — the schedule of an Edition (activities with days/times); staff are assigned to programs to mark attendance, and attendance grants certificates.
- **Edition check-in** — a door-level check-in record `{edition, user, checked_in_by, timestamp}`; made by any event staff (owner/admin/staff), repeatable, and a single check-in suffices for the edition certificate.
- **Ticket type / Product / Product variant** — what attendees can buy at an Edition.
- **Signature request** — a token-authenticated request to sign a document; completed requests produce a Signature.
- **Certification / Certification template** — a certificate design and its per-attendee emissions (by edition or by program).
- **Badge** — a per-person emission carrying the holder's Univents profile URL as a QR (`/profile/{id}`, the badge's *action URL*); networking + on-site attendance check-in. Emissions are keyed `(edition, holder, origin participant|staff)` — one badge per slot, updated in place, never a history of designs. The design is not stored on the emission: it is derived at read time from the edition's `badge_templates` (a ticket-type-specific template wins over the edition default; no match renders a placeholder), so design changes re-render live and nothing re-emails. Participant badges are emitted when a registration confirms (one email, with the QR embedded); staff badges are awarded on member add (published current/future editions only) and on edition publish, and past editions' staff badges are kept forever as keepsakes. Revoked emissions (`active|revoked` + reason) hide from profiles. The Univents backend never reads IdentityX profile data — holder names are derived by the front via IdentityX's public profile endpoint.
- **Intent** — a Payssage payment attempt (checkout); statuses flow from the provider.
- **Purchase** — an Edition-scoped order (tickets/products/program spots) created at checkout; items are reserved pending and confirmed or expired by a river worker; the record of truth for a buyer's shopping.
- **Doação (Donation)** — a gift to an Event (not a purchase): amount chosen by the donor (suggestions R$ 1/5/10 or free, minimum R$ 1); paid via a Payssage intent on the event's seller/wallet; recorded in its own table; no reservation or stock involved.
- **Wallet / Collector / Seller** — Payssage ownership units; payment methods (e.g. MercadoPago) connect to a Collector or Seller via OAuth.
- **Webhook endpoint / Webhook event / Webhook delivery** — Payssage fan-out of provider events to tenants.
- **Form / Namespace / Response / Step / Field** — Informd's multi-tenant form survey model.

## Architecture

- **Backend** — one of the four Go services (IdentityX, Payssage, Informd, Univents); a process, not a module.
- **Feature** — a slice of a Backend with its own models and routes (e.g. the `intents` feature).
- **TS API client** — the orval-generated TypeScript client per backend (functions + TanStack Query hooks + spec types) in `lib/ts/<svc>/client/`, produced by `just generate-orval` from each `api-spec.yml`; requests go through the shared fetcher stack via the `orval-mutator` (`@trieoh/api-client`), which unwraps the `fun.Response` envelope and rejects with `ApiError`.
- **Harness** — `lib/httpserver`: the shared HTTP-serving module every Backend boots through (server lifecycle, pprof, fun config, standard middleware stack, router skeleton, `/health`, `/metrics`, OpenTelemetry wrapping). One interface, four call sites. See ADR-0001.
- **Adapter** — a concrete thing that satisfies an interface at a seam (e.g. the MercadoPago provider module, the SQL repo modules, the email client).
- **Profile data boundary** — the Univents backend never reads IdentityX profile data (names, avatars, overrides); frontends derive it via IdentityX's public profile endpoint (`getActorProfile`). Univents emits only its own data: registration/ticket names, event/edition names, and the action URL.
