-- name: CreateUserSession :one
INSERT INTO user_sessions (token_id, user_id, issued_at, user_agent, user_ip, expires_at, project_id, user_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7,
        CASE
            WHEN $7::UUID IS NULL THEN 'client'
            ELSE 'project'
        END,
        NOW(), NOW())
RETURNING *;

-- name: GetUserSessionById :one
SELECT * FROM user_sessions
WHERE session_id = $1;

-- name: GetUserSessionByTokenId :one
SELECT * FROM user_sessions
WHERE token_id = $1;

-- name: ListUserSessions :many
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET
    issued_at = $2,
    user_agent = $3,
    user_ip = $4,
    expires_at = $5,
    token_id = $6,
    updated_at = NOW()
WHERE session_id = $1
RETURNING *;

-- name: DeleteUserSessionById :exec
DELETE FROM user_sessions
WHERE session_id = $1;

-- name: DeleteUserSessionByTokenId :exec
DELETE FROM user_sessions
WHERE token_id = $1;

-- name: RevokeUserSessionById :one
DELETE FROM user_sessions u
WHERE session_id = $1 AND token_id != $2 AND user_id = $3
RETURNING u.*;

-- name: RevokeOtherSessions :many
DELETE FROM user_sessions u
WHERE token_id != $1 AND user_id = $2
RETURNING u.*;

-- name: RevokeAllSessions :many
DELETE FROM user_sessions u
WHERE user_id = $1
RETURNING u.*;

-- name: DeleteExpiredSessions :many
DELETE FROM user_sessions u
WHERE expires_at < NOW()
RETURNING u.*;
