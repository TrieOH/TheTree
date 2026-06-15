-- name: HasActiveCryptoKey :one
SELECT EXISTS(
    SELECT 1
    FROM crypto_keys
    WHERE type = @type
      AND status = 'active'
      AND (project_id = @project_id OR (project_id IS NULL AND @project_id IS NULL))
);

-- name: CreateCryptoKey :one
INSERT INTO crypto_keys (project_id, type, status, public_key, encrypted_private_key, algorithm, expires_at) VALUES (
    @project_id,
    @type,
    'active',
    @public_key,
    @encrypted_private_key,
    @algorithm,
    @expires_at
) RETURNING *;

-- name: GetActiveCryptoKey :one
SELECT *
FROM crypto_keys
WHERE type = @type
  AND status = 'active'
  AND (project_id = @project_id OR (project_id IS NULL AND @project_id IS NULL));

-- name: GetCryptoKeyByID :one
SELECT *
FROM crypto_keys
WHERE id = @id;

-- name: GetActiveSigningKeys :many
SELECT id, public_key, algorithm
FROM crypto_keys
WHERE (project_id = @project_id OR @project_id IS NULL)
  AND type = 'signing'
  AND status IN ('active', 'retiring');