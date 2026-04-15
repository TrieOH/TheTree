-- name: CreateUserSession :one
INSERT INTO sessions (identity_id, issued_at, user_agent, user_ip, expires_at, project_id, user_type, created_at, updated_at)
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
SELECT s.*
FROM sessions s
JOIN identities i ON i.id = s.identity_id
WHERE i.type = $1
    AND i.entity_id = $2
    AND s.revoked_at IS NULL
    AND s.expires_at > NOW()
ORDER BY s.created_at DESC;

-- name: UpdateSession :exec
UPDATE sessions s
SET
    issued_at  = $4,
    user_agent = $5,
    user_ip    = $6,
    expires_at = $7,
    token_id   = $8,
    updated_at = NOW()
    FROM identities i
WHERE s.session_id = $1
  AND s.identity_id = i.id
  AND i.type = $2
  AND i.entity_id = $3
  AND s.revoked_at IS NULL;

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
UPDATE sessions s
SET
    revoked_at = NOW(),
    updated_at = NOW()
    FROM identities i
WHERE s.session_id = $1
  AND s.identity_id = i.id
  AND i.type = $2
  AND i.entity_id = $3
  AND s.revoked_at IS NULL
    RETURNING s.*;

-- name: RevokeSessionByFamilyID :exec
UPDATE sessions s
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE family_id = $1
AND s.revoked_at IS NULL;

-- name: RevokeOtherSessions :many
UPDATE sessions s
SET
    revoked_at = NOW(),
    updated_at = NOW()
FROM identities i
WHERE s.identity_id = i.id
AND i.type = $1
AND i.entity_id = $2
AND s.session_id != $3
AND s.revoked_at IS NULL
RETURNING s.*;

-- name: RevokeAllSessions :many
UPDATE sessions s
SET
    revoked_at = NOW(),
    updated_at = NOW()
FROM identities i
WHERE s.identity_id = i.id
AND i.type = $1
AND i.entity_id = $2
AND s.revoked_at IS NULL
RETURNING s.*;

-- name: RevokeExpiredSessions :many
UPDATE sessions
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE expires_at < NOW()
AND revoked_at IS NULL
RETURNING *;