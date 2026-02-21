-- name: GoAuthKeyExists :one
SELECT 1
FROM key_pair
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active'
    LIMIT 1;

-- name: CreateKeyPair :one
INSERT INTO key_pair (kid, project_id, key_type, algorithm, public_key, private_key, usage, status, expires_at, verify_expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetActiveSigningKeyForGoAuth :one
SELECT *
FROM key_pair
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active'
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
    LIMIT 1;

-- name: RotateSigningKeysForGoAuth :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active';

-- name: ListActivePublicKeysForGoAuth :many
SELECT
    kid,
    algorithm,
    public_key,
    created_at,
    expires_at
FROM key_pair
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND status IN ('active', 'rotated')
  AND verify_expires_at > now()
ORDER BY created_at DESC;

-- name: GetActiveSigningKeyForProject :one
SELECT *
FROM key_pair
WHERE
    project_id = $1
  AND key_type = 'project'
  AND usage = 'sign'
  AND status = 'active'
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
    LIMIT 1;

-- name: GetGoAuthKeyByKID :one
SELECT *
FROM key_pair
WHERE
    kid = $1
  AND key_type = 'goauth'
  AND status != 'revoked';

-- name: GetProjectKeyByKID :one
SELECT *
FROM key_pair
WHERE
    kid = $1
  AND key_type = 'project'
  AND status != 'revoked';

-- name: RotateSigningKeysForProject :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    project_id = $1
  AND key_type = 'project'
  AND usage = 'sign'
  AND status = 'active';

-- name: RevokeExpiredRotatedKeys :exec
UPDATE key_pair
SET status = 'revoked'
WHERE
    status = 'rotated'
  AND verify_expires_at < now();

-- name: RevokeKeyByKID :exec
UPDATE key_pair
SET status = 'revoked'
WHERE kid = $1;

-- name: ListActivePublicKeysForProject :many
SELECT
    kid,
    algorithm,
    public_key,
    created_at,
    expires_at
FROM key_pair
WHERE
    project_id = $1
  AND status IN ('active', 'rotated')
  AND verify_expires_at > now()
ORDER BY created_at DESC;

-- name: DeleteExpiredRevokedKeys :exec
DELETE FROM key_pair
WHERE
    status = 'revoked'
  AND expires_at < now();

-- name: GetActiveGoAuthSigningKID :one
SELECT kid
FROM key_pair
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active'
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
    LIMIT 1;

-- name: GetActiveProjectSigningKID :one
SELECT kid
FROM key_pair
WHERE
    project_id = $1
  AND key_type = 'project'
  AND usage = 'sign'
  AND status = 'active'
  AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
    LIMIT 1;

-- name: RotateExpiredGoAuthKeys :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active'
  AND expires_at < now();

-- name: RotateExpiredProjectKeys :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    key_type = 'project'
  AND usage = 'sign'
  AND status = 'active'
  AND expires_at < now();

-- INFRA ONLY --

-- name: ListProjectsWithActiveSigningKeys :many
SELECT DISTINCT project_id
FROM key_pair
WHERE
    key_type = 'project'
  AND usage = 'sign'
  AND status = 'active'
  AND project_id IS NOT NULL;

