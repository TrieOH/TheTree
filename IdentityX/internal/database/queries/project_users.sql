-- name: RegisterProjectUser :one
INSERT INTO project_users (project_id, email, password_hash, metadata)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProjectUserById :one
SELECT pu.*
FROM project_users AS pu
JOIN projects AS p on p.id = pu.project_id
WHERE pu.id = $1 AND pu.project_id = $2 AND p.owner_id = $3;

-- name: GetProjectUserByIdInternal :one
SELECT pu.*
FROM project_users AS pu
         JOIN projects AS p on p.id = pu.project_id
WHERE pu.id = $1 AND pu.project_id = $2;

-- name: ListProjectUsersExternal :many
SELECT pu.*
FROM project_users AS pu
JOIN projects AS p on p.id = pu.project_id
WHERE pu.project_id = $1
  AND p.owner_id = $2
ORDER BY pu.created_at DESC;

-- name: ListProjectUsersInternal :many
SELECT * FROM project_users
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: GetProjectUserByEmailExternal :one
SELECT pu.*
FROM project_users AS pu
JOIN projects AS p ON p.id = pu.project_id
WHERE pu.project_id = $1
  AND pu.email = $2
  AND p.owner_id = $3;

-- name: GetProjectUserByEmailInternal :one
SELECT * FROM project_users
WHERE project_id = $1 AND email = $2;

-- name: UpdateProjectUser :one
UPDATE project_users AS pu
SET
    email = $3,
    password_hash = $4,
    updated_at = NOW()
FROM projects AS p
WHERE pu.id = $1
  AND pu.project_id = $2
  AND p.id = pu.project_id
  AND p.owner_id = $5
    RETURNING pu.*;

-- name: DeleteProjectUser :exec
DELETE FROM project_users AS pu
USING projects AS p
WHERE pu.id = $1
  AND pu.project_id = $2
  AND p.id = pu.project_id
  AND p.owner_id = $3;

-- name: VerifyProjectUser :one
WITH updated AS (
UPDATE project_users AS u
SET
    is_verified = TRUE,
    verified_at = NOW()
WHERE u.id = sqlc.arg(user_id)
  AND u.is_verified = FALSE
    RETURNING TRUE
)
SELECT
    COALESCE(
            (SELECT TRUE FROM updated),
            FALSE
    )::boolean;

-- name: ProjectUserBelongsToProject :one
SELECT EXISTS (
    SELECT 1
    FROM project_users
    WHERE project_id = $1
     AND id = $2
);