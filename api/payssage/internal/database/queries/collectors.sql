-- name: CreateCollector :one
INSERT INTO collectors (provider, provider_user_id, credentials)
VALUES (@provider, @provider_user_id, @credentials)
RETURNING *;

-- name: GetCollectorByID :one
SELECT *
FROM collectors
WHERE id = @id;

-- name: ListCollectors :many
SELECT *
FROM collectors;
