-- name: BlacklistToken :exec
INSERT INTO refresh_blacklist (token_id, access_jti, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW());

-- name: GetRefreshBlacklistById :one
SELECT * FROM refresh_blacklist
WHERE token_id = $1;

-- name: GetRefreshBlacklistByAccessJTI :one
SELECT * FROM refresh_blacklist
WHERE access_jti = $1;

-- name: DeleteRefreshBlacklist :exec
DELETE FROM refresh_blacklist
WHERE token_id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM refresh_blacklist
WHERE expires_at < NOW();
