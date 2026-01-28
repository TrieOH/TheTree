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

-- name: PermissionBelongsToProject :one
SELECT EXISTS(
    SELECT 1
    FROM permissions
    WHERE id = $1
        AND project_id = $2
);

----------------------------------------------
------------- Direct Permissions -------------
----------------------------------------------

-- name: GiveDirectPermission :exec
INSERT INTO identity_permissions (identity_id, permission_id, scope_id)
VALUES ($1, $2, $3);

-- name: TakeDirectPermission :exec
DELETE FROM identity_permissions
WHERE identity_id = $1 AND permission_id = $2 AND scope_id IS NOT DISTINCT FROM $3;

-----------------------------------------------
------------ Permission Resolution ------------
-----------------------------------------------

-- name: GetEffectivePermissions :many
WITH effective_ids AS (
    -- Direct permission assignments (scope-specific + inherited global WITHIN same realm)
    SELECT ip.permission_id AS id
    FROM identity_permissions ip
    WHERE ip.identity_id = $1
        AND (
            ip.scope_id IS NOT DISTINCT FROM $2
                OR ($2 IS NOT NULL AND ip.scope_id IS NULL)
            )

    UNION

    -- Role-based permission assignments (scope-specific + inherited global WITHIN same realm)
    SELECT rp.permission_id AS id
    FROM identity_roles ir
    JOIN role_permissions rp ON rp.role_id = ir.role_id
    WHERE ir.identity_id = $1
        AND (
            ir.scope_id IS NOT DISTINCT FROM $2
                OR ($2 IS NOT NULL AND ir.scope_id IS NULL)
            )
)
SELECT p.*
FROM permissions p
JOIN effective_ids e ON p.id = e.id
WHERE
    -- STRICT isolation: IdP and Project permissions never mix
    (sqlc.narg(project_id)::UUID IS NULL AND p.project_id IS NULL)      -- IdP realm only
    OR (sqlc.narg(project_id)::UUID IS NOT NULL AND p.project_id = sqlc.narg(project_id)::UUID)  -- Project realm only
ORDER BY p.object, p.action, p.id;