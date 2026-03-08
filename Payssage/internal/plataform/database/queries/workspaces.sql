-- name: CreateWorkspace :one
INSERT INTO workspaces (id, scope_id, user_id, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWorkspaceByID :one
SELECT *
FROM workspaces
WHERE id = $1;

-- name: GetWorkspaceByName :one
SELECT *
FROM workspaces
WHERE user_id = $1 AND name = $2;

-- name: ListWorkspacesByUser :many
SELECT *
FROM workspaces
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: EnableSandbox :one
UPDATE workspaces
SET sandbox = true, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DisableSandbox :one
UPDATE workspaces
SET sandbox = false, updated_at = now()
WHERE id = $1
RETURNING *;