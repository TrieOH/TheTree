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

-- name: GetActiveSigningKey :one
SELECT *
FROM key_pair
WHERE
    project_id IS NOT DISTINCT FROM @project_id
    AND key_type = CASE WHEN @project_id::uuid IS NULL THEN 'goauth' ELSE 'project' END
  AND usage = 'sign'
  AND status = 'active'
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: ListPublicKeys :many
SELECT
    kid,
    algorithm,
    public_key,
    created_at,
    expires_at
FROM key_pair
WHERE
    project_id IS NOT DISTINCT FROM @project_id
    AND key_type = CASE WHEN @project_id::uuid IS NULL THEN 'goauth' ELSE 'project' END
  AND status IN ('active', 'rotated')
  AND verify_expires_at > now()
ORDER BY created_at DESC;

-- name: GetKeyByKID :one
SELECT *
FROM key_pair
WHERE kid = $1
  AND status != 'revoked';

-- name: RotateSigningKeys :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    project_id IS NOT DISTINCT FROM @project_id
    AND key_type = CASE WHEN @project_id::uuid IS NULL THEN 'goauth' ELSE 'project' END
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

-- name: DeleteExpiredRevokedKeys :exec
DELETE FROM key_pair
WHERE
    status = 'revoked'
  AND expires_at < now();

-- name: GetActiveSigningKID :one
SELECT kid
FROM key_pair
WHERE
    project_id IS NOT DISTINCT FROM @project_id
    AND key_type = CASE WHEN @project_id::uuid IS NULL THEN 'goauth' ELSE 'project' END
  AND usage = 'sign'
  AND status = 'active'
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: RotateExpiredKeys :exec
UPDATE key_pair
SET
    status = 'rotated',
    usage = 'verify'
WHERE
    usage = 'sign'
  AND status = 'active'
  AND expires_at < now();

-- INFRA ONLY --

-- name: ListProjectsWithSigningKeys :many
SELECT DISTINCT project_id
FROM key_pair
WHERE
    key_type = 'project'
  AND project_id IS NOT NULL;

