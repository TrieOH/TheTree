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

-- name: AddNamespaceMember :exec
INSERT INTO namespace_members (user_id, namespace_id, role, added_at, added_by)
VALUES ($1, $2, $3, $4, $5);

-- name: RemoveNamespaceMember :exec
DELETE FROM namespace_members
WHERE user_id = $1 AND namespace_id = $2;

-- name: GetNamespaceMember :one
SELECT *
FROM namespace_members
WHERE user_id = $1 AND namespace_id = $2;

-- name: ListNamespaceMembers :many
SELECT *
FROM namespace_members
WHERE namespace_id = $1;

-- name: ListOwnedNamespaces :many
SELECT *
FROM namespaces
WHERE owner_id = $1;

-- name: ListJoinedNamespaces :many
SELECT ns.*
FROM namespace_members nm
INNER JOIN namespaces ns
ON nm.namespace_id = ns.id
WHERE nm.user_id = $1;