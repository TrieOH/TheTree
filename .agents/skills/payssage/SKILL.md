---
name: payssage
description: Use the Payssage service — wallets, sellers, payment intents, tenant webhooks, test-mode simulation. Use when working on or calling api/payssage, creating/configuring the platform wallet (sandbox, fee, webhook endpoint), creating or cancelling intents, wiring PAYSSAGE_* envs, signing/verifying webhook deliveries, or testing payment flows (including /testmode).
---

# Payssage

Payssage is the payments service: **wallets** (ownership units), **sellers/collectors** (provider accounts like MercadoPago), **intents** (payment attempts), and a **webhook pipeline** (provider events → HMAC-signed deliveries to tenant endpoints). Univents runs on **one platform wallet** (env-configured, boot-verified) — see `docs/plans/issue-61-checkout.md` D6/D2/D3 for the store design.

Run locally: `docker compose up payssage` → `http://localhost:8082` (env in `api/payssage/.env`). Contract: `api/payssage/api-spec.yml`.

## Key concepts

| Concept | Meaning |
|---|---|
| Wallet | ownership unit; has `owner_id`, optional `organization_id`, `sandbox`, `fee_bps`. Intents/webhooks/sellers hang off it |
| Seller | provider account bound to a wallet (splits transactions); **created only via the provider OAuth callback** — no create-seller API |
| Collector | wallet's payment collector (receives funds); bound via OAuth |
| Intent | one payment attempt; statuses `pending | processing | succeeded | cancelled | failed | rejected | refunded` |
| Webhook endpoint | per-wallet `{name, url}` → returns a **secret** used to verify `X-Payssage-Signature` on deliveries |
| Webhook delivery | envelope body (below), HMAC-SHA256(secret, raw body), retried up to 5× on non-2xx |

## Auth: two kinds of callers

- **User JWT** — a user in the payssage identityx project (bearer). Most wallet/intent/org ops are bearerAuth-only.
- **Service API key** — the payssage svc-account key (`X-API-Key`); accepted on wallet create/get/fee and (after this repo's spec fixes) sandbox, webhook endpoints, and testmode — so the platform can manage the wallet with the service key alone.

**Wallet ownership matters**: `CheckWalletAccess` requires the caller to be the wallet owner (or an org member of an org-scoped wallet). Univents' boot check calls `GetWallet` with the payssage service key, so **the platform wallet must be owned by the payssage svc actor** (create it with the payssage service key, not a user JWT).

## Bootstrap the platform wallet (order matters)

```bash
PKEY=<payssage svc key>   # X-API-Key
# 1. create (owner = payssage svc actor)
curl -s -X POST /wallets -H "X-API-Key: $PKEY" -d '{"name":"Platform Wallet"}'   # → $WALLET_ID
# 2. sandbox (PATCH, not POST!)
curl -s -X PATCH /wallets/$WALLET_ID/sandbox -H "X-API-Key: $PKEY" -d '{"sandbox":true}'
# 3. marketplace fee (bps; 500 = 5%)
curl -s -X PATCH /wallets/$WALLET_ID/fee -H "X-API-Key: $PKEY" -d '{"fee_bps":500}'
# 4. webhook endpoint → returns the secret (store as PAYSSAGE_WEBHOOK_SECRET)
curl -s -X POST /wallets/$WALLET_ID/webhooks/endpoints -H "X-API-Key: $PKEY" \
  -d '{"name":"univents-store","url":"http://univents:8080/webhooks/payssage"}'
```

Wire the consuming service's env: `PAYSSAGE_WALLET_ID`, `PAYSSAGE_WEBHOOK_SECRET`, `PAYSSAGE_API_KEY` (= the **payssage** service key — it is the wallet owner).

## Intents

- `POST /wallets/{wallet_id}/checkout` — real checkout (payload: payment_method, amount_cents, seller, etc.); pix returns the QR in `provider_data`, cards charge synchronously.
- `POST /intents/{intent_id}/cancel`, `GET /intents/{intent_id}`.
- **Test mode** (`TEST_MODE=true` in the payssage env — otherwise `/testmode/*` 503s): `POST /testmode/intents/create` hard-creates an intent with a chosen `status` (e.g. `succeeded`) so downstream flows observe the webhook without a real provider. Callable with the service key.

## Webhook contract (D2 envelope — Univents consumes this)

Payssage delivers, per endpoint, a body that is NOT the raw provider payload:

```json
{ "intent_id": "...", "wallet_id": "...", "provider": "mercadopago",
  "external_id": "<provider_payment_id>", "event_type": "payment.succeeded",
  "payload": { "...": "raw provider payload" } }
```

`X-Payssage-Signature` = hex(HMAC-SHA256(raw body, endpoint secret)). Verify against the **exact bytes** POSTed. Return 200 to stop retries; non-2xx retries (up to 5 ≈ 5s).

## Gotchas (learned the hard way)

- `setWalletSandbox` is **PATCH** — POST returns 405.
- The webhook URL must be reachable from payssage's container (docker network hostnames, not `localhost`).
- The 000-owner bug: an API-key identity must come from **unwrapping the introspect envelope** (`data.subject`); reading the envelope itself zeroes the identity and every wallet created that way gets `owner_id = 00000000-…` (fixed in this repo's `internal/app/auth.go` — keep the unwrap).
- Sellers only exist after a real provider OAuth (`/providers/mercado_pago/connect` → browser → callback). For local e2e you can seed the `sellers` table directly (wallet_id, provider, provider_user_id, credentials JSON with the MP test access token/public key).
- `listWalletSellers` verifies ownership; the seller used by univents must belong to the platform wallet or `completeEventPayments` 403s.
- Dropping the payssage DB wipes wallets/webhooks/intents — recreate the wallet and re-point `PAYSSAGE_WALLET_ID`/`PAYSSAGE_WEBHOOK_SECRET`.

## Pointers

- `api/payssage/api-spec.yml` — wallets, sellers, intents, webhooks, testmode, oauth.
- `docs/plans/issue-61-*.md` — the store's payssage integration (envelope D2, one-wallet D6, webhook-only approval D3, refund deferral D1).
- Sibling skills: `identityx` (the payssage project + svc key), `univents` (webhook receiver, checkout, boot check).
