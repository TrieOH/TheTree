-- name: CreateIntent :one
INSERT INTO intents (id, wallet_id, seller_id, collector_id, amount_cents, currency, sandbox, provider, status, provider_data, metadata)
VALUES (@id, @wallet_id, @seller_id, @collector_id, @amount_cents, @currency, @sandbox, @provider, @status, @provider_data, @metadata)
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

-- name: CancelIntent :one
UPDATE intents
SET status = 'cancelled', updated_at = NOW()
WHERE id = @id
  AND status = 'processing'
RETURNING *;

-- name: ConfirmIntent :one
UPDATE intents
SET status = 'succeeded', updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: FailIntent :one
UPDATE intents
SET status = 'failed', updated_at = NOW()
WHERE id = @id
RETURNING *;

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
