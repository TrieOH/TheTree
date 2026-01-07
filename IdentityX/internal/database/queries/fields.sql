-- name: CreateSchemaField :one
INSERT INTO schema_fields (
    id,
    schema_id,
    schema_version_id,
    key,
    type,
    owner,
    title,
    description,
    placeholder,
    required,
    mutable,
    default_value,
    position
)
SELECT
    sqlc.arg(field_id),
    sv.schema_id,
    sv.id,
    sqlc.arg(key),
    sqlc.arg(type),
    sqlc.arg(owner),
    sqlc.arg(title),
    sqlc.arg(description),
    sqlc.arg(placeholder),
    sqlc.arg(required),
    sqlc.arg(mutable),
    sqlc.arg(default_value),
    sqlc.arg(position)
FROM schema_versions sv
WHERE
    sv.id = sqlc.arg(schema_version_id)
  AND sv.schema_id = sqlc.arg(schema_id)
  AND sv.status = 'draft'
    RETURNING *;

-- name: GetFieldsByVersionID :many
SELECT *
FROM schema_fields
WHERE schema_version_id = $1;