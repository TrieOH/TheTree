-- name: CreateForm :one
INSERT INTO forms (namespace_id, owner_id, name, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFormByID :one
SELECT * FROM forms
WHERE id = $1;

-- name: BulkGetForms :many
SELECT * FROM forms
WHERE id = ANY($1::uuid[])
ORDER BY created_at DESC;