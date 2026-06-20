-- name: CreateApiKey :one
INSERT INTO api_keys (actor_id, project_id, name, key_prefix, key_hash, expires_at)
VALUES (@actor_id, @project_id, @name, @key_prefix, @key_hash, @expires_at)
RETURNING *;

-- name: ListApiKeysByProject :many
SELECT *
FROM api_keys
WHERE project_id = @project_id;

-- name: GetApiKeyByPrefix :one
SELECT *
FROM api_keys
WHERE key_prefix = @key_prefix;

-- name: GetApiKeyByPrefixAndProject :one
SELECT *
FROM api_keys
WHERE project_id = @project_id
  AND key_prefix = @key_prefix;

-- name: RevokeApiKeyByPrefixAndProject :exec
UPDATE api_keys
SET
    revoked_at = NOW()
WHERE project_id = @project_id
  AND key_prefix = @key_prefix;

-- name: SetApiKeyLastUsedAtByPrefixAndProject :exec
UPDATE api_keys
SET
    last_used_at = NOW()
WHERE project_id = @project_id
  AND key_prefix = @key_prefix;
