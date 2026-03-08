-- name: CreateIntent :one
INSERT INTO intents (id, workspace_id, amount, currency, status, client_secret, provider, metadata)
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
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

-- name: PayIntent :one
UPDATE intents
SET
    status = $2,
    provider_payment_id = $3,
    updated_at = now()
WHERE id = $1 AND status = 'pending'
    RETURNING *;

-- name: GetIntentByProviderPaymentID :one
SELECT * FROM intents WHERE provider_payment_id = $1;