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
WITH target_scope AS (
    SELECT
        s.id,
        s.project_id,
        s.name,
        s.external_id
    FROM scopes s
    WHERE s.id = sqlc.narg('scope_id')::uuid
    AND sqlc.narg('scope_id')::uuid IS NOT NULL
),
master_scope AS (
    SELECT s2.id
    FROM target_scope ts
        JOIN scopes s2 ON s2.project_id = ts.project_id
        AND s2.name = ts.name
        AND s2.external_id IS NULL
    WHERE ts.external_id IS NOT NULL
),
scope_family AS (
    -- The specific scope itself
    SELECT id FROM target_scope
    UNION
    -- Its parent master scope (same name, no external_id)
    SELECT id FROM master_scope
)
SELECT DISTINCT p.*
FROM permissions p
WHERE p.id IN (
    -- Role-derived permissions
    SELECT rp.permission_id
    FROM identity_roles ir
    JOIN role_permissions rp ON rp.role_id = ir.role_id
    WHERE ir.identity_id = sqlc.arg('identity_id')::uuid
    AND (
        -- Querying nil scope: only nil scope assignments
        (sqlc.narg('scope_id')::uuid IS NULL AND ir.scope_id IS NULL)
        OR
        -- Querying specific scope: hierarchy + nil scope fallback
        (sqlc.narg('scope_id')::uuid IS NOT NULL AND (
            ir.scope_id IN (SELECT id FROM scope_family)
            OR ir.scope_id IS NULL
        ))
    )

    UNION

    -- Direct permission assignments
    SELECT ip.permission_id
    FROM identity_permissions ip
    WHERE ip.identity_id = sqlc.arg('identity_id')::uuid
    AND (
        (sqlc.narg('scope_id')::uuid IS NULL AND ip.scope_id IS NULL)
        OR
        (sqlc.narg('scope_id')::uuid IS NOT NULL AND (
            ip.scope_id IN (SELECT id FROM scope_family)
            OR ip.scope_id IS NULL
        ))
    )
)
AND (
    -- Strict IdP/Project isolation
    (sqlc.narg('project_id')::uuid IS NULL AND p.project_id IS NULL)
    OR (sqlc.narg('project_id')::uuid IS NOT NULL AND p.project_id = sqlc.narg('project_id')::uuid)
)
ORDER BY p.object, p.action;