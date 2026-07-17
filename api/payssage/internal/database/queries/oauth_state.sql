-- name: CreateOAuthState :one
INSERT INTO oauth_states (state, wallet_id, provider, flow, final_redirect_url, expires_at)
VALUES (@state, @wallet_id, @provider, @flow, @final_redirect_url, @expires_at)
RETURNING *;

-- name: GetOAuthState :one
SELECT *
FROM oauth_states
WHERE state = @state;

-- name: DeleteOAuthState :exec
DELETE FROM oauth_states
WHERE state = @state;