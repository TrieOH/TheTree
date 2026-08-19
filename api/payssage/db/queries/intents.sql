-- name: CreateIntent :one
INSERT INTO intents (id, wallet_id, seller_id, collector_id, amount_cents, currency, sandbox, provider, status, status_detail, provider_data, metadata, external_id, external_group)
VALUES (@id, @wallet_id, @seller_id, @collector_id, @amount_cents, @currency, @sandbox, @provider, @status, @status_detail, @provider_data, @metadata, @external_id, @external_group)
RETURNING *;

-- name: GetIntentByID :one
SELECT *
FROM intents
WHERE id = @id;

-- name: ListIntentsByWallet :many
SELECT *
FROM intents
WHERE wallet_id = @wallet_id
ORDER BY created_at DESC;

-- name: UpdateIntentProviderData :one
UPDATE intents
SET provider_data = @provider_data, updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: GetIntentBySellerID :one
SELECT *
FROM intents
WHERE seller_id = @seller_id
LIMIT 1;

-- name: GetIntentByCollectorID :one
SELECT *
FROM intents
WHERE collector_id = @collector_id
LIMIT 1;

-- name: ListIntentsByOwner :many
SELECT i.*
FROM intents i
JOIN wallets w ON i.wallet_id = w.id
WHERE w.owner_id = @owner_id
  AND w.organization_id IS NULL
ORDER BY i.created_at DESC;

-- name: ListIntentsByOrg :many
SELECT i.*
FROM intents i
JOIN wallets w ON i.wallet_id = w.id
WHERE w.organization_id = @organization_id
ORDER BY i.created_at DESC;

-- name: UpdateIntent :one
UPDATE intents
SET status = @status,
    status_detail = @status_detail,
    provider_data = @provider_data,
    updated_at = NOW()
WHERE id = @id
    RETURNING *;

-- name: GetIntentByProviderTransactionID :one
SELECT *
FROM intents
WHERE provider = @provider
  AND provider_data->>'transaction_id' = @transaction_id::text;