-- name: CreateUserSession :one
INSERT INTO user_sessions (token_id, issued_at, user_agent, user_ip, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING *;

-- name: GetUserSessionById :one
SELECT * FROM user_sessions
WHERE session_id = $1;

-- name: GetUserSessionByTokenId :one
SELECT * FROM user_sessions
WHERE token_id = $1;

-- name: ListUserSessions :many
SELECT * FROM user_sessions ORDER BY created_at DESC;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET
    issued_at = $2,
    user_agent = $3,
    user_ip = $4,
    expires_at = $5,
    updated_at = NOW()
WHERE session_id = $1
RETURNING *;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions
WHERE session_id = $1;

-- name: DeleteUserSessionByTokenId :exec
DELETE FROM user_sessions
WHERE token_id = $1;

-- name: RevokeUserSession :one
DELETE FROM user_sessions u
WHERE session_id = $1 AND token_id != $2
RETURNING u.token_id;
