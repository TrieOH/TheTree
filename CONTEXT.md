# TrieOH — Domain & Architecture Context

Single-context glossary for TheTree. Architecture decisions live in `docs/adr/`.
Terms here give names to good seams; keep this current when concepts are added
or sharpened.

## Domain

- **Actor** — a human, service, or machine identity in IdentityX; everything that can authenticate.
- **Organization** — top-level tenant in IdentityX; members hold roles (member/admin/owner).
- **Project** — an IdentityX workspace under an organization; scope for API keys and profiles.
- **Event** — a Univents gathering; owned by members with roles (owner/admin/staff).
- **Edition** — one run of an Event, with its own dates and registration window.
- **Program / Program occurrence** — the schedule of an Edition (activities with days/times).
- **Ticket type / Product / Product variant** — what attendees can buy at an Edition.
- **Signature request** — a token-authenticated request to sign a document; completed requests produce a Signature.
- **Certification / Certification template** — a certificate design and its per-attendee emissions (by edition or by program).
- **Intent** — a Payssage payment attempt (checkout); statuses flow from the provider.
- **Wallet / Collector / Seller** — Payssage ownership units; payment methods (e.g. MercadoPago) connect to a Collector or Seller via OAuth.
- **Webhook endpoint / Webhook event / Webhook delivery** — Payssage fan-out of provider events to tenants.
- **Form / Namespace / Response / Step / Field** — Informd's multi-tenant form survey model.

## Architecture

- **Backend** — one of the four Go services (IdentityX, Payssage, Informd, Univents); a process, not a module.
- **Feature** — a slice of a Backend with its own models and routes (e.g. the `intents` feature).
- **Harness** — `lib/httpserver`: the shared HTTP-serving module every Backend boots through (server lifecycle, pprof, fun config, standard middleware stack, router skeleton, `/health`, `/metrics`, OpenTelemetry wrapping). One interface, four call sites. See ADR-0001.
- **Adapter** — a concrete thing that satisfies an interface at a seam (e.g. the MercadoPago provider module, the SQL repo modules, the email client).
