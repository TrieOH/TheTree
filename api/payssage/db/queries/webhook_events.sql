-- name: CreateWebhookEvent :one
INSERT INTO webhook_events (wallet_id, intent_id, provider, external_id, event_type, status_detail, payload)
VALUES (@wallet_id, @intent_id, @provider, @external_id, @event_type, @status_detail, @payload)
    RETURNING *;

-- name: GetWebhookEventByID :one
SELECT *
FROM webhook_events
WHERE id = @id;

-- name: ListWebhookEventsByWallet :many
SELECT *
FROM webhook_events
WHERE wallet_id = @wallet_id;

-- name: ListWebhookEventsByProvider :many
SELECT *
FROM webhook_events
WHERE provider = @provider;
