-- name: CreateOAuthState :one
INSERT INTO oauth_states (state, workspace_id, provider, flow, is_marketplace, fee_bps, final_redirect_url, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *;

-- name: GetOAuthState :one
SELECT * FROM oauth_states
WHERE state = $1 AND expires_at > now();

-- name: DeleteOAuthState :exec
DELETE FROM oauth_states WHERE state = $1;

-- name: CreateProviderCredential :one
INSERT INTO provider_credentials (workspace_id, provider, credentials)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: GetProviderCredential :one
SELECT * FROM provider_credentials
WHERE id = $1 AND revoked_at IS NULL;

-- name: ListProviderCredentials :many
SELECT * FROM provider_credentials
WHERE workspace_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeProviderCredential :one
UPDATE provider_credentials
SET revoked_at = now()
WHERE id = $1 AND workspace_id = $2 AND revoked_at IS NULL
    RETURNING *;

-- name: GetWorkspaceProviderCredential :one
SELECT * FROM provider_credentials
WHERE workspace_id = $1
  AND provider = $2
  AND revoked_at IS NULL
ORDER BY created_at DESC
    LIMIT 1;

-- FIXME Make this query take in credential ID for the seller
-- name: GetSellerCredentialByProvider :one
SELECT pc.* FROM provider_credentials pc
 LEFT JOIN marketplace_configs mc ON mc.credential_id = pc.id
WHERE pc.workspace_id = $1
  AND pc.provider = $2
  AND pc.revoked_at IS NULL
  AND mc.id IS NULL
ORDER BY pc.created_at DESC
    LIMIT 1;