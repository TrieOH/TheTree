-- name: CreateApiKey :one
INSERT INTO api_keys (subject_id, name, display_prefix, key_hash, created_by, expires_at)
VALUES (@subject_id, @name, @display_prefix, @key_hash, @created_by, @expires_at)
RETURNING *;

-- name: GetActorApiKeys :many
SELECT *
FROM api_keys
WHERE subject_id = @subject_id;

-- name: GetApiKeyByPrefix :one
SELECT *
FROM api_keys
WHERE display_prefix = @display_prefix
  AND revoked_at IS NULL;

-- name: SetApiKeyLastUsedAtByPrefix :exec
UPDATE api_keys
SET
    last_used_at = NOW()
WHERE display_prefix = @display_prefix;