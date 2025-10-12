-- name: CreateUsers :one
INSERT INTO users (email, password, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;

-- name: GetUsersById :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUsers :one
UPDATE users
SET
    email = $2,
    password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE id = $1;
