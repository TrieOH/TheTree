-- name: CreatePermission :one
INSERT INTO permissions (object, action, conditions, project_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPermissionByIDInternal :one
SELECT *
FROM permissions
WHERE id = $1;

-- name: GetPermissionByIDExternal :one
SELECT *
FROM permissions
WHERE id = $1 AND project_id = $2;

-- name: ListPermissionsByProject :many
SELECT *
FROM permissions
WHERE
    -- project boundary (required)
    project_id = sqlc.arg(project_id)

    -- optional object filter
    AND (
        sqlc.narg(object)::text IS NULL
        OR object = sqlc.narg(object)
    )

    -- optional action filter
    AND (
        sqlc.narg(action)::text IS NULL
        OR action = sqlc.narg(action)
    )
ORDER BY created_at DESC;
