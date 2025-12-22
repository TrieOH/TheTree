-- name: RevokeToken :exec
INSERT INTO refresh_blacklist (token_id, expires_at, created_at)
VALUES ($1, $2, NOW());

-- name: RevokeManyTokens :many
INSERT INTO refresh_blacklist (token_id, expires_at)
SELECT UNNEST($1::uuid[]), UNNEST($2::timestamp[])
ON CONFLICT (token_id) DO NOTHING
RETURNING token_id;

-- name: GetRevokedRefreshByID :one
SELECT * FROM refresh_blacklist
WHERE token_id = $1;

-- name: DeleteRevokedRefreshByID :exec
DELETE FROM refresh_blacklist
WHERE token_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_blacklist
WHERE expires_at < NOW();
