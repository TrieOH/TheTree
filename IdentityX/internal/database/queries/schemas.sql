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

-- name: SchemaExists :one
SELECT EXISTS (
    SELECT 1
    FROM schemas
    WHERE project_id = $1 AND flow_id = $2 AND type = $3
);

-- name: SchemaBelongsToProject :one
SELECT EXISTS (
    SELECT 1
    FROM schemas
    WHERE id = $1
      AND project_id = $2
);

-- name: GetSchema :one
SELECT *
FROM schemas
WHERE id = $1 AND project_id = $2;

-- name: GetSchemaByFlowID :one
SELECT *
FROM schemas
WHERE flow_id = $1 AND project_id = $2;

-- name: ListSchemas :many
SELECT *
FROM schemas
WHERE project_id = $1;

-- name: SetSchemaVersion :exec
UPDATE schemas
SET
    current_version_id = $1
WHERE id = $2 AND project_id = $3;

-- /////////////////////////////// --
-- //// -- Schema Versions -- //// --
-- /////////////////////////////// --

-- name: DraftSchemaVersion :one
INSERT INTO schema_versions (schema_id, version)
VALUES ($1, $2)
RETURNING *;

-- name: PublishSchemaVersion :exec
UPDATE schema_versions s
SET
    status = 'published',
    updated_at = NOW()
WHERE id = $1
  AND schema_id = $2
  AND status = 'draft';

-- name: ArchiveSchemaVersion :exec
UPDATE schema_versions
SET
    status = 'archived',
    updated_at = NOW()
WHERE id = $1 AND schema_id = $2;

-- name: GetLatestSchemaVersion :one
SELECT *
FROM schema_versions
WHERE schema_id = $1
ORDER BY version DESC
LIMIT 1;
