-- name: CreateSessionIdentity :one
INSERT INTO identities (type, entity_id, created_at)
VALUES ($1, $2, NOW())
    RETURNING *;

-- name: GetSessionIdentityByEntityIDAndType :one
SELECT *
FROM identities
WHERE entity_id = $1 AND type = $2;

-- name: GetSessionIdentityByIDAndType :one
SELECT *
FROM identities
WHERE id = $1 AND type = $2;