-- name: CreateMarketplaceConfig :one
INSERT INTO marketplace_configs (workspace_id, credential_id, fee_bps, provider)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: ListMarketplaceConfigs :many
SELECT * FROM marketplace_configs
WHERE workspace_id = $1
ORDER BY created_at;

-- name: GetMarketplaceConfig :one
SELECT * FROM marketplace_configs
WHERE workspace_id = $1 AND credential_id = $2;

-- name: UpdateMarketplaceConfig :one
UPDATE marketplace_configs
SET
    fee_bps    = $3,
    updated_at = now()
WHERE workspace_id = $1 AND credential_id = $2
    RETURNING *;

-- name: DeleteMarketplaceConfig :exec
DELETE FROM marketplace_configs
WHERE workspace_id = $1 AND credential_id = $2;

-- name: DeleteAllMarketplaceConfigs :exec
DELETE FROM marketplace_configs
WHERE workspace_id = $1;

-- name: GetMarketplaceConfigByProvider :one
SELECT mc.* FROM marketplace_configs mc
JOIN provider_credentials pc ON pc.id = mc.credential_id
WHERE mc.workspace_id = $1 AND pc.provider = $2;