-- name: CreateSchemaField :one
INSERT INTO schema_fields (
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
WHERE schema_version_id = $1
ORDER BY position;

-- name: ListFieldsFromSchema :many
SELECT *
FROM schema_fields
WHERE schema_id = $1
ORDER BY schema_version_id, position;

-- name: CloneFields :execrows
INSERT INTO schema_fields (
    object_id,
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
    position,
    created_at,
    updated_at
)
SELECT
    gen_random_uuid(),
    f.id,
    dv.schema_id,
    dv.id,
    f.key,
    f.type,
    f.owner,
    f.title,
    f.description,
    f.placeholder,
    f.required,
    f.mutable,
    f.default_value,
    f.position,
    NOW(),
    NOW()
FROM schema_fields f
JOIN schema_versions dv ON dv.id = sqlc.arg(draft_version_id)::UUID
WHERE f.schema_version_id = sqlc.arg(source_version_id)::UUID;

-- name: DiffVersionFields :one
WITH base_fields AS (
    SELECT
        id,
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
    FROM schema_fields
    WHERE schema_version_id = sqlc.arg(base_version_id)::UUID
),
draft_fields AS (
    SELECT
        id,
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
    FROM schema_fields
    WHERE schema_version_id = sqlc.arg(draft_version_id)::UUID
) SELECT EXISTS (
    -- base → draft differences
    (
        SELECT * FROM base_fields
        EXCEPT
        SELECT * FROM draft_fields
    )
    UNION ALL
    -- draft → base differences
    (
        SELECT * FROM draft_fields
        EXCEPT
        SELECT * FROM base_fields
    )
) AS has_changes;
