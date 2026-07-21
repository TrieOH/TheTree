-- name: CreateWebhookEndpoint :one
INSERT INTO webhook_endpoints (wallet_id, name, url, secret)
VALUES (@wallet_id, @name, @url, @secret)
    RETURNING *;

-- name: GetWebhookEndpointByID :one
SELECT *
FROM webhook_endpoints
WHERE id = @id;

-- name: ListWebhookEndpointsByWallet :many
SELECT *
FROM webhook_endpoints
WHERE wallet_id = @wallet_id;

-- name: DeleteWebhookEndpoint :exec
DELETE FROM webhook_endpoints
WHERE id = @id;
