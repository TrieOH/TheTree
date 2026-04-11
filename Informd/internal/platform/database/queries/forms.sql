-- name: CreateForm :one
INSERT INTO forms (id, project_id, owner_id, scope_id, title, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFormsByProject :many
SELECT * FROM forms
WHERE project_id = $1
ORDER BY created_at DESC;