-- name: UpsertApiKey :exec
INSERT INTO api_keys (project_id, client_id, key_hash, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (project_id) DO UPDATE
SET key_hash = EXCLUDED.key_hash,
    client_id = EXCLUDED.client_id,
    updated_at = NOW();

-- name: GetApiKeyByProjectID :one
SELECT project_id, client_id, key_hash, created_at, updated_at
FROM api_keys
WHERE project_id = $1;

-- name: DeleteApiKey :exec
DELETE FROM api_keys
WHERE project_id = $1;
