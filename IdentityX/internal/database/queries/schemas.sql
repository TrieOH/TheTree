-- name: DraftSchema :one
INSERT INTO schemas (project_id, title, flow_id, type, status)
VALUES ($1, $2, $3, $4, 'draft')
RETURNING *;

-- name: PublishSchema :exec
UPDATE schemas
SET
    status = 'published',
    updated_at = NOW()
WHERE id = $1 AND project_id = $2 AND status = 'draft';

-- name: ArchiveSchema :exec
UPDATE schemas
SET
    status = 'archived',
    updated_at = NOW()
WHERE id = $1 AND project_id = $2 AND status = 'published';

-- name: DeleteSchema :exec
DELETE FROM schemas
WHERE id = $1 AND project_id = $2 AND status = 'archived';

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

-- name: GetSchemaByFlowIDAndType :one
SELECT *
FROM schemas
WHERE flow_id = $1 AND type = $2 AND project_id = $3;

-- name: ListSchemas :many
SELECT *
FROM schemas
WHERE project_id = $1;

-- name: SetSchemaVersion :exec
UPDATE schemas
SET
    current_version_id = $1
WHERE id = $2 AND project_id = $3 AND status = 'published';

-- /////////////////////////////// --
-- //// -- Schema Versions -- //// --
-- /////////////////////////////// --

-- name: DraftSchemaVersion :one
INSERT INTO schema_versions (schema_id, version)
VALUES ($1, $2)
RETURNING *;

-- name: PublishSchemaVersion :execrows
UPDATE schema_versions
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
WHERE id = $1 AND schema_id = $2 AND status = 'published';

-- name: GetLatestSchemaVersion :one
SELECT *
FROM schema_versions
WHERE schema_id = $1
ORDER BY version DESC
LIMIT 1;

-- name: GetLatestSchemaVersionForUpdate :one
SELECT *
FROM schema_versions
WHERE schema_id = $1
ORDER BY version DESC
    LIMIT 1
FOR UPDATE;

-- name: GetCurrentSchemaVersion :one
SELECT sv.*
FROM schemas s
JOIN schema_versions sv
ON sv.id = s.current_version_id
WHERE s.id = $1;

-- name: ListSchemaVersion :many
SELECT *
FROM schema_versions
WHERE schema_id = $1
ORDER BY version DESC;

-- name: CopyVersionOnDraft :one
WITH locked AS (
    SELECT *
    FROM schema_versions sv
    WHERE sv.id = $1
      AND status = 'published'
    FOR UPDATE
    )
INSERT INTO schema_versions (schema_id, version, status, based_on_version_id)
SELECT
    locked.schema_id,
    locked.version + 1,
    'draft',
    locked.id
FROM locked
RETURNING *;


