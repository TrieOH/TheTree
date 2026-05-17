-- name: CreateUserSession :one
INSERT INTO sessions (user_id, issued_at, user_agent, user_ip, expires_at, project_id, user_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6,
        CASE
            WHEN $6::UUID IS NULL THEN 'client'
            ELSE 'project'
        END,
        NOW(), NOW())
RETURNING *;

-- name: GetUserSessionByID :one
SELECT *
FROM sessions
WHERE session_id = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: GetSessionByFamilyID :one
SELECT *
FROM sessions
WHERE family_id = $1;

-- name: GetUserSessionByTokenID :one
SELECT *
FROM sessions
WHERE token_id = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: ListSessions :many
SELECT *
FROM sessions
WHERE user_type = $1
    AND user_id = $2
    AND revoked_at IS NULL
    AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: UpdateSession :exec
UPDATE sessions
SET
    issued_at  = $4,
    user_agent = $5,
    user_ip    = $6,
    expires_at = $7,
    token_id   = $8,
    updated_at = NOW()
WHERE session_id = $1
  AND user_id = $2
  AND user_type = $3
  AND entity_id = $4
  AND revoked_at IS NULL;

-- name: DeleteRevokedSessions :many
DELETE FROM sessions
WHERE revoked_at IS NOT NULL
    RETURNING *;

-- name: RotateSessionToken :one
UPDATE sessions
SET
    expires_at = $1,
    token_id   = sqlc.arg(new_token_id)::UUID,
    issued_at  = NOW(),
    updated_at = NOW()
WHERE family_id = $2
  AND token_id = sqlc.arg(old_token_id)::UUID
  AND revoked_at IS NULL
  AND expires_at > NOW()
    RETURNING *;

-- ============================
-- Revocation (soft, auditable)
-- ============================

-- name: RevokeSessionByID :one
UPDATE sessions
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE session_id = $1
  AND user_id = $2
  AND user_type = $3
  AND revoked_at IS NULL
    RETURNING *;

-- name: RevokeSessionByFamilyID :exec
UPDATE sessions
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE family_id = $1
AND revoked_at IS NULL;

-- name: RevokeOtherSessions :many
UPDATE sessions
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE user_id = $1
  AND user_type = $2
  AND session_id != $3
  AND revoked_at IS NULL
    RETURNING *;

-- name: RevokeAllSessions :many
UPDATE sessions s
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE user_id = $1
  AND user_type = $2
  AND revoked_at IS NULL
    RETURNING *;

-- name: RevokeExpiredSessions :many
UPDATE sessions
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE expires_at < NOW()
  AND revoked_at IS NULL
    RETURNING *;