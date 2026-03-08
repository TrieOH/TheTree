-- name: CreateMarketplaceConfig :one
INSERT INTO marketplace_configs (workspace_id, credential_id, fee_bps)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: GetMarketplaceConfig :one
SELECT * FROM marketplace_configs
WHERE workspace_id = $1;

-- name: UpdateMarketplaceConfig :one
UPDATE marketplace_configs
SET
    credential_id = $2,
    fee_bps       = $3,
    updated_at    = now()
WHERE workspace_id = $1
    RETURNING *;

-- name: DeleteMarketplaceConfig :exec
DELETE FROM marketplace_configs
WHERE workspace_id = $1;