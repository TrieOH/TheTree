-- name: CreateRole :one
INSERT INTO roles (name, description, project_id)
VALUES ($1, $2, $3)
RETURNING *;

-- pass nil to project_id if you want GoAuth roles
-- name: UpdateRoleDescription :exec
UPDATE roles
SET
    description = $1,
    updated_at = NOW()
WHERE id = $2 AND project_id = $3;

-- name: GetRoleByIDInternal :one
SELECT *
FROM roles
WHERE id = $1;

-- name: GetRoleByIDExternal :one
SELECT *
FROM roles
WHERE id = $1 AND project_id = $2;

-- pass nil to project_id if you want GoAuth roles
-- name: GetRoleByName :one
SELECT *
FROM roles
WHERE name = $1 AND project_id = $2;

-- name: ListRolesByProject :many
SELECT *
FROM roles
WHERE project_id = $1
ORDER BY created_at DESC;