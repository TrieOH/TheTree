-- name: CreateProject :one
INSERT INTO projects (
  project_name, owner_id, metadata, is_active,
  pub_key, priv_key
)
VALUES (
  $1, $2, $3, $4, $5,
  pgp_sym_encrypt(sqlc.arg(priv_key)::text, current_setting('app.jwt_master_key'))
)
RETURNING
  id, project_name, owner_id, metadata, is_active, pub_key,
  pgp_sym_decrypt(priv_key::bytea, current_setting('app.jwt_master_key')) AS priv_key,
  created_at, updated_at;

-- name: GetProjectById :one
SELECT 
  id, project_name, owner_id, metadata, is_active,
  pub_key, created_at, updated_at
FROM projects
WHERE id = $1 AND owner_id = $2;

-- name: GetProjectKeysById :one
SELECT pub_key, pgp_sym_decrypt(priv_key::bytea, current_setting('app.jwt_master_key')) AS priv_key
FROM projects
WHERE id = $1 AND owner_id = $2;

-- name: GetProjectPublicKeyById :one
SELECT pub_key
FROM projects
WHERE id = $1;

-- name: ListProjects :many
SELECT 
  id, project_name, owner_id, metadata, is_active, created_at, updated_at
FROM projects
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET 
  project_name = $3, metadata = $4,
  updated_at = NOW()
WHERE id = $1 and owner_id = $2
RETURNING
    id, project_name, owner_id, metadata, is_active, pub_key,
    created_at, updated_at;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1 AND owner_id = $2;

-- admin
-- name: AdminGetProjectById :one
SELECT 
  id, project_name, owner_id, metadata, is_active,
  created_at, updated_at
FROM projects
WHERE id = $1;


