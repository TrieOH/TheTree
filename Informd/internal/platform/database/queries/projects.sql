-- name: CreateProject :one
INSERT INTO projects (id, owner_id, name)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: GetProjectByID :one
SELECT *
FROM projects
WHERE id = $1;

-- name: GetProjectByName :one
SELECT *
FROM projects
WHERE owner_id = $1 AND name = $2;

-- name: ListProjectsByOwner :many
SELECT *
FROM projects
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: ListProjectsByIDs :many
SELECT *
FROM projects
WHERE id = ANY($1::uuid[])
ORDER BY created_at DESC;