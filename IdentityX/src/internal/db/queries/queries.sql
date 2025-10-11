-- name: CreateUser :one
INSERT INTO users (first_name, last_name, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users 
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
