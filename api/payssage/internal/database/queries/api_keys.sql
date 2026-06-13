-- name: CreateAPIKey :one
INSERT INTO api_keys (id, scope_id, workspace_id, name, key_hash, key_prefix)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAPIKeyByPrefix :many
SELECT *
FROM api_keys
WHERE key_prefix = $1 AND revoked_at IS NULL;

-- name: ListAPIKeysByWorkspace :many
SELECT *
FROM api_keys
WHERE workspace_id = $1
ORDER BY created_at DESC;

-- name: RevokeAPIKey :one
UPDATE api_keys
SET
    revoked_at = now()
WHERE id = $1
  AND workspace_id = $2
  AND revoked_at IS NULL
RETURNING *;