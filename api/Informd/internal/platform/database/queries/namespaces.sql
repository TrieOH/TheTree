-- name: CreateNamespace :one
INSERT INTO namespaces (owner_id, name)
VALUES ($1, $2)
    RETURNING *;

-- name: GetNamespaceByID :one
SELECT *
FROM namespaces
WHERE id = $1;

-- name: GetNamespaceByName :one
SELECT *
FROM namespaces
WHERE owner_id = $1 AND name = $2;

-- name: BulkGetNamespaces :many
SELECT * FROM namespaces
WHERE id = ANY($1::uuid[])
ORDER BY created_at DESC;