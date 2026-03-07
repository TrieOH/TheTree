-- name: CreateWebhookEndpoint :one
INSERT INTO webhook_endpoints (id, scope_id, workspace_id, url, secret)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWebhookEndpointByID :one
SELECT * FROM webhook_endpoints
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListWebhookEndpointsByWorkspace :many
SELECT * FROM webhook_endpoints
WHERE workspace_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteWebhookEndpoint :exec
UPDATE webhook_endpoints
SET deleted_at = now()
WHERE id = $1 AND workspace_id = $2;

-- name: CreateWebhookDelivery :one
INSERT INTO webhook_deliveries (id, endpoint_id, intent_id, event, payload, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetWebhookDeliveryByID :one
SELECT *
FROM webhook_deliveries
WHERE id = $1;

-- name: ListWebhookDeliveriesByEndpoint :many
SELECT *
FROM webhook_deliveries
WHERE endpoint_id = $1
ORDER BY created_at DESC;

-- name: MarkDeliveryDelivered :one
UPDATE webhook_deliveries
SET
    status = 'delivered',
    attempts = attempts + 1,
    last_attempted_at = now()
WHERE id = $1
RETURNING *;

-- name: MarkDeliveryFailed :one
UPDATE webhook_deliveries
SET
    status = 'failed',
    attempts = attempts + 1,
    last_attempted_at = now()
WHERE id = $1
RETURNING *;

-- name: IncrementDeliveryAttempt :one
UPDATE webhook_deliveries
SET
    attempts = attempts + 1,
    last_attempted_at = now()
WHERE id = $1
RETURNING *;