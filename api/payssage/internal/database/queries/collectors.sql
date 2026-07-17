-- name: CreateCollector :one
INSERT INTO collectors (owner_id, organization_id, provider, provider_user_id, credentials)
VALUES (@owner_id, @organization_id, @provider, @provider_user_id, @credentials)
RETURNING *;

-- name: GetCollectorByID :one
SELECT *
FROM collectors
WHERE id = @id;

-- name: ListCollectors :many
SELECT *
FROM collectors;

-- name: ListCollectorsByOwner :many
SELECT *
FROM collectors
WHERE owner_id = @owner_id
  AND organization_id IS NULL;

-- name: ListCollectorsByOrg :many
SELECT *
FROM collectors
WHERE organization_id = @organization_id;
