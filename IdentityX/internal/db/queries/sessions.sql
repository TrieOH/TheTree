-- name: CreateUserSession :one
INSERT INTO sessions (user_id, issued_at, user_agent, user_ip, expires_at, project_id, user_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6,
        CASE
            WHEN $6::UUID IS NULL THEN 'client'
            ELSE 'project'
        END,
        NOW(), NOW())
RETURNING *;

-- name: GetUserSessionById :one
SELECT * FROM sessions
WHERE session_id = $1;

-- name: GetUserSessionByTokenId :one
SELECT * FROM sessions
WHERE token_id = $1;

-- name: ListUserSessions :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateUserSession :exec
UPDATE sessions
SET
    issued_at = $2,
    user_agent = $3,
    user_ip = $4,
    expires_at = $5,
    token_id = $6,
    updated_at = NOW()
WHERE session_id = $1;

-- name: DeleteSessionsByFilter :many
DELETE FROM sessions
WHERE
    user_id = $1
  AND (sqlc.narg(session_id)::uuid IS NULL OR session_id = sqlc.narg(session_id))
  AND (sqlc.narg(exclude_id)::uuid IS NULL OR session_id != sqlc.narg(exclude_id))
  AND (sqlc.narg(token_id)::uuid IS NULL OR token_id = sqlc.narg(token_id))
  AND (sqlc.narg(expired_before)::timestamp IS NULL OR expires_at < sqlc.narg(expired_before))
    RETURNING *;

-- name: DeleteExpiredSessions :many
DELETE FROM sessions u
WHERE expires_at < NOW()
RETURNING u.*;
