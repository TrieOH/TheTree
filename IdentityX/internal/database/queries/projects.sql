-- name: CreateProject :one
INSERT INTO projects (project_name, owner_id, metadata, is_active, domain)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProjectByIDExternal :one
SELECT *
FROM projects
WHERE id = $1 AND owner_id = $2;

-- name: GetProjectByIDInternal :one
SELECT *
FROM projects
WHERE id = $1;

-- name: ListProjects :many
SELECT *
FROM projects
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET 
  project_name = $3,
  metadata = $4,
  domain = $5,
  updated_at = NOW()
WHERE id = $1 and owner_id = $2
RETURNING *;

-- name: DeleteProject :execrows
DELETE FROM projects
WHERE id = $1 AND owner_id = $2;

-- name: IsOwnerOf :one
SELECT EXISTS (
    SELECT 1
    FROM projects
    WHERE id = $1 AND owner_id = $2
);


