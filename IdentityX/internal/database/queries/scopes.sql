-- name: CreateScope :one
INSERT INTO scopes (type, project_id, name, external_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

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
