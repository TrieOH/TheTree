---
name: univents
description: Use the Univents service — events, editions, ticketing, products, programs, badges, certifications, and the store (purchases, webhooks, realtime). Use when working on or calling api/univents, building/seeding a catalog (event → edition → tickets/products/programs → badges/certs), running a purchase end-to-end (webhook receiver, read surfaces, WS/SSE), wiring UNIVENTS_* / PAYSSAGE_* envs, or following the issue-61 store splits.
---

# Univents

Univents is the events service: events & editions, ticket types, products, programs (schedules), registrations, badge & certification templates, signatures, and — via the issue-61 store — **purchases** with payssage payment intents, a webhook receiver, read surfaces, and realtime (WS + SSE).

Run locally: `docker compose up univents` → `http://localhost:8081` (env in `api/univents/.env`). Contract: `api/univents/api-spec.yml`. The store design master is `docs/plans/issue-61-checkout.md`; each split doc (`issue-61-split-N.md`) is its scope.

## Key concepts

| Concept | Meaning |
|---|---|
| Event / Edition | event (draft/active) owns editions (runs with dates + registration window); editions publish |
| Ticket type / Product variant / Program occurrence | the three purchasable item kinds (`item_type`: `ticket`\|`product`\|`program_occurrence`) |
| Registration / Product purchase / Participation | the materialized rows a purchase owns (pending → confirmed/cancelled/expired) |
| Purchase | the store's record of truth: `purchases` + `purchase_items` (availability ledger) + `ws_tokens`; statuses `pending|approved|expired|cancelled`; correlation by `payssage_intent_id` |
| Badge / Cert | templates with a `design_data` JSON (below); participant badges emit on registration confirm |

## Auth

Bearer JWT from the **univents** identityx project. Public (anonymous) reads: event/edition browsing, products, programs, certification verification, and the SSE store stream. Owner-only ops return **404** for non-owners (no existence leak).

## Bootstrap a catalog (curl recipes, verified)

```bash
JWT=<univents project user jwt>   # e.g. admin@trieoh.com / Admin123#
# event → payments → edition → publish
POST /events                          {"full_name","acronym","slug","description","contact_email"}          → $EVENT_ID
POST /events/$EVENT_ID/payments/complete   {"seller_id","public_key"}   # seller must exist on the platform wallet
POST /events/$EVENT_ID/editions      {"name","slug","starts_at","ends_at"}                                 → $EDITION_ID
POST /events/$EVENT_ID/editions/$EDITION_ID/publish
# catalog
POST /editions/$EDITION_ID/ticket-types    {"name","price_cents","access_level","max_quantity"}
POST /editions/$EDITION_ID/products        {"vendor_code","variant_vendor_code","name","price","stock"}    # product + first variant
PATCH /variants/$VARIANT_ID                {"vendor_code","name","gallery_urls":[...]}                     # web images OK
POST /editions/$EDITION_ID/programs        {"kind":"activity","name","price","banner_url","min_access_level"}
POST /programs/$PROGRAM_ID/occurrences     {"starts_at","ends_at","max_capacity"}
POST /editions/$EDITION_ID/badges          {"name","origin":null|"staff","design_data":{...}}
POST /editions/$EDITION_ID/certifications/templates   {"name","kind":"edition_attendance"|"program_attendance","design_data":{...}}
POST /certifications/templates/$TEMPLATE_ID/link      {"program_id"}   # program-attendance kind
```

**design_data shapes** (mirror the frontend editors):
- Badge: `{"canvas":{"width":321,"height":204},"backgroundColor":"#…","background":"<img url or null>","elements":[{"type":"text","x","y","width","height","paragraphs":[{"align","lineHeight","runs":[{"text","fontSize","fontFamily","color","bold","italic","underline"}]}]} | {"type":"image","src","fit","radius","opacity"} | {"type":"qr","value":"{{checkin_url}}","foreground","background","style"}]}`
- Cert: canvas ~1123×794; elements `text` (same rich text), `image` (`src`,`fit`,`radius`,`opacity`), `hash` (`hashLabel`,`hash":"{{cert_hash}}"`,`linkLabel`,`url":"{{verify_url}}"`), `signature` (`signatureId`,`src`,`name`). Text variables: `{{participant_name}} {{event_name}} {{edition_name}} {{activity_name}} {{location}} {{certified_at}}`.

## The store: buy flow (issue-61 splits)

Checkout `POST /editions/{id}/checkout` is **split 7 — not landed**. Until then the buy path is exercised on seeded data (splits 4–6):

1. Seed a purchase (mirror split-7's write): `purchases` (pending, edition, purchaser, total, `payssage_intent_id`, expires_at) + `purchase_items` + the materialized `registrations` row in one tx.
2. Create the intent in payssage: `POST /testmode/intents/create` (status `succeeded`, sandbox) — needs payssage `TEST_MODE=true`.
3. Fire the webhook: `POST /webhooks/payssage` with the D2 envelope + `X-Payssage-Signature: hex(HMAC-SHA256(PAYSSAGE_WEBHOOK_SECRET, raw body))`. → 200 → purchase `approved`, registration `confirmed`, badge emitted, `NOTIFY univents_changes`.
4. Verify surfaces: `GET /checkouts/{id}` (owner-only), `GET /purchases`, `GET /ws/token?purchase_id=…` → `WS /ws?token=…` (snapshot on connect, closes on terminal), `GET /editions/{id}/store/stream` (SSE snapshot + `event: stock` deltas).

**Realtime contract**: WS frames are `purchase.snapshot` / `intent.updated` / `purchase.confirmed` / `purchase.expired` / `purchase.cancelled` (`{"type","payload"}`); tokens are 32-byte random, SHA-256 at rest, 10-min TTL, one-time. SSE items are `{"id","item_type","stock"}` with `stock: null` = unlimited; numbers always recomputed from the DB (publishers send only item ids).

## Gotchas (learned the hard way)

- Checkout isn't merged — don't try `POST /editions/{id}/checkout` yet; the buy e2e is seed + signed webhook.
- `completeEventPayments` verifies the seller belongs to the platform wallet (split 2); a seeded seller must match `PAYSSAGE_WALLET_ID`'s wallet.
- **`patchProductVariant` does not merge**: omitting `stock` sets it NULL. Include every field you intend to keep when patching (open issue).
- SSE is a **raw hijacked route** — the harness has a 60s WriteTimeout and no Flush through its middleware; the stream hijacks the connection. Don't "fix" it by buffering.
- The WS socket is per-purchase and one-time-token; reconnect = fresh `GET /ws/token` + re-open (the snapshot restores context).
- Boot fails fast on the platform wallet: `PAYSSAGE_WALLET_ID` must resolve via payssage `GetWallet` (wallet owned by the payssage svc actor) or univents panics.
- Images: direct web URLs work anywhere the frontend renders them; uploads go through rustfs (`OBJECT_STORAGE_*` in `api/univents/.env` — must match rustfs creds, `USE_SSL=false` locally).
- The constraint registry (`internal/app/constraints.go`) must name every constraint the DB actually has or boot fails — when editing migrations in place (dev mode), keep it in sync.

## Pointers

- `api/univents/api-spec.yml` — every op; raw WS/SSE routes are documented in the info description block.
- `docs/plans/issue-61-checkout.md` + `issue-61-split-{1..7}.md` — the store design, decisions (D1–D10), and per-split scope.
- `CONTEXT.md` — domain terms (Event, Edition, Badge, Purchase, Wallet…).
- Sibling skills: `payssage` (intents, webhook envelope, wallet), `identityx` (users, project, service key).
