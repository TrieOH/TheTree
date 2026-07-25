-- name: CreateWebhookDelivery :one
INSERT INTO webhook_deliveries (endpoint_id, event_id, status)
VALUES (@endpoint_id, @event_id, @status)
    RETURNING *;

-- name: GetWebhookDeliveryByID :one
SELECT *
FROM webhook_deliveries
WHERE id = @id;

-- name: ListWebhookDeliveriesByEndpoint :many
SELECT *
FROM webhook_deliveries
WHERE endpoint_id = @endpoint_id;

-- name: MarkDeliveryDelivered :one
UPDATE webhook_deliveries
SET
    status = 'delivered',
    last_attempted_at = NOW(),
    attempts = attempts + 1
WHERE id = @id
    RETURNING *;

-- name: MarkDeliveryFailed :one
UPDATE webhook_deliveries
SET
    status = 'failed',
    last_attempted_at = NOW(),
    attempts = attempts + 1
WHERE id = @id
    RETURNING *;

-- name: IncrementDeliveryAttempt :one
UPDATE webhook_deliveries
SET
    attempts = attempts + 1,
    last_attempted_at = NOW()
WHERE id = @id
    RETURNING *;

-- name: UpdateWebhookDelivery :one
UPDATE webhook_deliveries
SET status            = @status,
    attempts          = @attempts,
    last_attempted_at = @last_attempted_at,
    response_status   = @response_status,
    response_body     = @response_body,
    updated_at        = NOW()
WHERE id = @id
    RETURNING *;