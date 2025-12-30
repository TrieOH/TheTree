-- name: RevokeToken :exec
INSERT INTO revoked_refresh_tokens (token_id, expires_at, created_at)
VALUES ($1, $2, NOW());

-- name: RevokeManyTokens :many
INSERT INTO revoked_refresh_tokens (token_id, expires_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::timestamp[])
ON CONFLICT (token_id) DO NOTHING
RETURNING token_id;

-- name: GetRevokedRefreshByID :one
SELECT * FROM revoked_refresh_tokens
WHERE token_id = $1;

-- name: DeleteRevokedRefreshByID :exec
DELETE FROM revoked_refresh_tokens
WHERE token_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM revoked_refresh_tokens
WHERE expires_at < NOW();

-- name: IsRefreshTokenRevoked :one
SELECT EXISTS (
    SELECT 1
    FROM revoked_refresh_tokens
    WHERE token_id = $1
);