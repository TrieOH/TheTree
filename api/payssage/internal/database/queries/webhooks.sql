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

-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (id, provider, event_type, payload)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: EnrichWebhookEvent :one
UPDATE webhook_events
SET
    workspace_id = $2,
    intent_id    = $3,
    external_id  = $4
WHERE id = $1
    RETURNING *;

-- name: GetWebhookEventByID :one
SELECT * FROM webhook_events
WHERE id = $1;

-- name: ListWebhookEventsByWorkspace :many
SELECT * FROM webhook_events
WHERE workspace_id = $1
ORDER BY received_at DESC;

-- name: ListWebhookEventsByProvider :many
SELECT * FROM webhook_events
WHERE provider = $1
ORDER BY received_at DESC;