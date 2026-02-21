-- name: CreateRole :one
INSERT INTO roles (name, description, project_id, meta)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- pass nil to project_id if you want GoAuth roles
-- name: UpdateRoleDescription :exec
UPDATE roles
SET
    description = $1,
    updated_at = NOW()
WHERE id = $2 AND project_id = $3;

-- pass nil to project_id if you want GoAuth roles
-- name: UpdateRoleMeta :exec
UPDATE roles
SET
    meta = $1,
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

-- name: RoleBelongsToProject :one
SELECT EXISTS (
    SELECT 1
    FROM roles
    WHERE id = $1
      AND project_id = $2
);

------------------------------
------ Role Permissions ------
------------------------------

-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2);

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 and permission_id = $2;

-- name: GetRolePermissions :many
SELECT p.*
FROM permissions p
INNER JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1 AND p.project_id = $2
ORDER BY created_at DESC;

------------------------------------
------- Role User Management -------
------------------------------------

-- name: GiveRole :exec
INSERT INTO identity_roles (identity_id, role_id, scope_id)
VALUES ($1, $2, $3);

-- name: TakeRole :exec
DELETE FROM identity_roles
WHERE identity_id = $1 AND role_id = $2 AND scope_id IS NOT DISTINCT FROM $3;

-- name: GetUserRoles :many
SELECT r.*, ir.scope_id, s.name as scope_name, s.external_id
FROM roles r
INNER JOIN identity_roles ir ON ir.role_id = r.id
LEFT JOIN scopes s ON s.id = ir.scope_id
WHERE ir.identity_id = $1
  AND r.project_id = $2
ORDER BY
    r.name,
    ir.scope_id IS NULL DESC,  -- Project-wide roles first (NULL scope_id)
    s.name NULLS FIRST,        -- project_root scopes before named scopes (if applicable)
    s.external_id NULLS FIRST; -- Master scopes (NULL) before specific resources