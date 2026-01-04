-- name: DraftSchema :one
INSERT INTO schemas (project_id, title, flow_id, type, status)
VALUES ($1, $2, $3, $4, 'draft')
RETURNING *;

-- name: PublishSchema :exec
UPDATE schemas
SET
    status = 'published',
    updated_at = NOW()
WHERE id = $1 AND project_id = $2;

-- name: ArchiveSchema :exec
UPDATE schemas
SET
    status = 'archived',
    updated_at = NOW()
WHERE id = $1 AND project_id = $2;

-- name: DeleteSchema :exec
DELETE FROM schemas
WHERE id = $1 AND project_id = $2;

-- name: GetSchema :one
SELECT *
FROM schemas
WHERE id = $1 AND project_id = $2;

-- name: ListSchemas :many
SELECT *
FROM schemas
WHERE project_id = $1;
