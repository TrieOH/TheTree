-- name: CreateIntent :one
INSERT INTO intents (id, workspace_id, amount, currency, status, provider, provider_data, metadata, seller_credential_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING *;

-- name: GetIntentByID :one
SELECT *
FROM intents
WHERE id = $1;

-- name: ListIntents :many
SELECT *
FROM intents
ORDER BY created_at DESC;

-- name: ListIntentsByWorkspace :many
SELECT *
FROM intents
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: CancelIntent :one
UPDATE intents
SET
    status = 'cancelled',
    updated_at = now()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: ConfirmIntent :one
UPDATE intents
SET
    status = 'succeeded',
    updated_at = now()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: FailIntent :one
UPDATE intents
SET
    status = 'failed',
    updated_at = now()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: UpdateIntentProviderData :one
UPDATE intents
SET
    provider_data = provider_data || $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetIntentByMPOrderID :one
SELECT *
FROM intents
WHERE provider_data->>'order_id' = $1::text
  AND provider = 'mercadopago';