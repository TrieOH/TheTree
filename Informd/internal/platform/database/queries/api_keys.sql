-- name: CreateAPIKey :one
INSERT INTO api_keys (owner_id, name, key_hash, key_prefix)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAPIKeyByPrefix :many
SELECT *
FROM api_keys
WHERE key_prefix = $1 AND revoked_at IS NULL;

-- name: ListAPIKeys :many
SELECT *
FROM api_keys
WHERE owner_id = $1 AND (revoked_at IS NULL OR revoked_at >= now() - INTERVAL '30 days')
ORDER BY created_at DESC;

-- name: RevokeAPIKey :one
UPDATE api_keys
SET
    revoked_at = now()
WHERE id = $1
  AND owner_id = $2
  AND revoked_at IS NULL
RETURNING *;