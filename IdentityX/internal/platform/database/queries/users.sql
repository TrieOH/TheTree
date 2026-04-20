-- name: RegisterUser :one
INSERT INTO users (email, password_hash, project_id, user_type)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ResetUserPassword :exec
UPDATE users
SET password_hash = $1
WHERE id = $2;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByEmailFromProject :one
SELECT * FROM users
WHERE email = $1
  AND project_id = $2;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    password_hash = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: VerifyUser :one
WITH updated AS (
UPDATE users AS u
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
    FROM users
    WHERE project_id = $1
      AND id = $2
);

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
    last_login_at = now()
WHERE id = $1;

-- name: ListUsersFromProject :many
SELECT * FROM users
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: GetUserByIDFromProject :one
SELECT * FROM users
WHERE id = $1
  AND project_id = $2;