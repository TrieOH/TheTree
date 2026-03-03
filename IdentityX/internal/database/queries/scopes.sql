-- name: CreateScope :one
INSERT INTO scopes (type, parent_id, project_id, name, external_id, meta)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProjectScopeMeta :exec
UPDATE scopes
SET
    meta = $1,
    updated_at = NOW()
WHERE id = $2 AND project_id = $3 AND type = 'project_scope';

-- name: DeleteScope :exec
DELETE FROM scopes
WHERE id = $1 AND project_id = $2 AND type = 'project_scope';

-- name: GetScopeByIDInternal :one
SELECT *
FROM scopes
WHERE id = $1;

-- name: GetScopeByIDExternal :one
SELECT *
FROM scopes
WHERE id = $1 AND project_id = $2;

-- internal only
-- name: GetRootByProjectID :one
SELECT *
FROM scopes
WHERE project_id = $1 AND type = 'project_root';

-- name: GetProjectScopes :many
SELECT *
FROM scopes
WHERE project_id = $1 AND type = 'project_scope';

-- name: GetGlobalScope :one
SELECT *
FROM scopes
WHERE type = 'global';

-- name: GetScopeAncestors :many
-- Gets all ancestors of a scope up to the root using recursive CTE
WITH RECURSIVE ancestors AS (
    -- Start with the target scope
    SELECT scopes.id, scopes.parent_id, scopes.type, scopes.project_id, scopes.name, scopes.external_id, scopes.created_at, scopes.meta, 0 as depth
    FROM scopes
    WHERE scopes.id = $1

    UNION ALL

    -- Recursively get all parents
    SELECT s.id, s.parent_id, s.type, s.project_id, s.name, s.external_id, s.created_at, s.meta, a.depth + 1
    FROM scopes s
    INNER JOIN ancestors a ON s.id = a.parent_id
    WHERE a.parent_id IS NOT NULL
)
SELECT ancestors.id, ancestors.parent_id, ancestors.type, ancestors.project_id, ancestors.name, ancestors.external_id, ancestors.created_at, ancestors.meta
FROM ancestors
ORDER BY ancestors.depth DESC;

-- name: GetScopeChildren :many
-- Gets direct children of a scope
SELECT * FROM scopes WHERE parent_id = $1;
