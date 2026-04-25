-- name: CreateForm :one
INSERT INTO forms (namespace_id, owner_id, name, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListFormsByUser :many
SELECT * FROM forms
WHERE owner_id = $1 AND namespace_id IS NULL
ORDER BY created_at DESC;

-- name: ListFormsByNamespace :many
SELECT * FROM forms
WHERE namespace_id = $1
ORDER BY created_at DESC;