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

-- name: CreateSchemaFieldsBatch :copyfrom
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
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12
 );

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

-- name: ListFieldsFromVersion :many
SELECT *
FROM schema_fields
WHERE schema_id = $1 AND schema_version_id = $2
ORDER BY position;

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
    uuidv7(),
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

-- name: DiffVersionFieldsFull :one
WITH
    base_fields AS (
        SELECT
            id, key, type, owner, title, description,
    placeholder, required, mutable, default_value, position
FROM schema_fields
WHERE schema_version_id = sqlc.arg(base_version_id)::UUID
    ),
    draft_fields AS (
SELECT
    id, key, type, owner, title, description,
    placeholder, required, mutable, default_value, position
FROM schema_fields
WHERE schema_version_id = sqlc.arg(draft_version_id)::UUID
    ),
    base_opts AS (
SELECT
    f.id as field_stable_id,
    opt.value, opt.label, opt.position
FROM schema_field_options opt
    JOIN schema_fields f ON f.object_id = opt.field_id
WHERE f.schema_version_id = sqlc.arg(base_version_id)::UUID
    ),
    draft_opts AS (
SELECT
    f.id as field_stable_id,
    opt.value, opt.label, opt.position
FROM schema_field_options opt
    JOIN schema_fields f ON f.object_id = opt.field_id
WHERE f.schema_version_id = sqlc.arg(draft_version_id)::UUID
    ),
    base_vis AS (
SELECT
    f.id as field_stable_id,
    dep_f.id as depends_on_stable_id,
    vr.operator, vr.value
FROM schema_field_visibility_rules vr
    JOIN schema_fields f ON f.object_id = vr.field_id
    JOIN schema_fields dep_f ON dep_f.object_id = vr.depends_on_field_id
WHERE f.schema_version_id = sqlc.arg(base_version_id)::UUID
  AND dep_f.schema_version_id = sqlc.arg(base_version_id)::UUID
    ),
    draft_vis AS (
SELECT
    f.id as field_stable_id,
    dep_f.id as depends_on_stable_id,
    vr.operator, vr.value
FROM schema_field_visibility_rules vr
    JOIN schema_fields f ON f.object_id = vr.field_id
    JOIN schema_fields dep_f ON dep_f.object_id = vr.depends_on_field_id
WHERE f.schema_version_id = sqlc.arg(draft_version_id)::UUID
  AND dep_f.schema_version_id = sqlc.arg(draft_version_id)::UUID
    ),
    base_req AS (
SELECT
    f.id as field_stable_id,
    dep_f.id as depends_on_stable_id,
    rr.operator, rr.value
FROM schema_field_required_rules rr
    JOIN schema_fields f ON f.object_id = rr.field_id
    JOIN schema_fields dep_f ON dep_f.object_id = rr.depends_on_field_id
WHERE f.schema_version_id = sqlc.arg(base_version_id)::UUID
  AND dep_f.schema_version_id = sqlc.arg(base_version_id)::UUID
    ),
    draft_req AS (
SELECT
    f.id as field_stable_id,
    dep_f.id as depends_on_stable_id,
    rr.operator, rr.value
FROM schema_field_required_rules rr
    JOIN schema_fields f ON f.object_id = rr.field_id
    JOIN schema_fields dep_f ON dep_f.object_id = rr.depends_on_field_id
WHERE f.schema_version_id = sqlc.arg(draft_version_id)::UUID
  AND dep_f.schema_version_id = sqlc.arg(draft_version_id)::UUID
    )
SELECT
    (EXISTS(SELECT 1 FROM (SELECT * FROM base_fields EXCEPT SELECT * FROM draft_fields) x)
        OR EXISTS(SELECT 1 FROM (SELECT * FROM draft_fields EXCEPT SELECT * FROM base_fields) x))
        as fields_changed,

    (EXISTS(SELECT 1 FROM (SELECT * FROM base_opts EXCEPT SELECT * FROM draft_opts) x)
        OR EXISTS(SELECT 1 FROM (SELECT * FROM draft_opts EXCEPT SELECT * FROM base_opts) x))
        as options_changed,

    (EXISTS(SELECT 1 FROM (SELECT * FROM base_vis EXCEPT SELECT * FROM draft_vis) x)
        OR EXISTS(SELECT 1 FROM (SELECT * FROM draft_vis EXCEPT SELECT * FROM base_vis) x))
        as visibility_rules_changed,

    (EXISTS(SELECT 1 FROM (SELECT * FROM base_req EXCEPT SELECT * FROM draft_req) x)
        OR EXISTS(SELECT 1 FROM (SELECT * FROM draft_req EXCEPT SELECT * FROM base_req) x))
        as required_rules_changed;

-- /////////////////////////////// --
-- //// -- FIELD OPTIONS -- //// --
-- /////////////////////////////// --

-- name: CreateFieldOption :one
INSERT INTO schema_field_options (
    field_id,
    value,
    label,
    position
)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: CreateFieldOptionsBatch :copyfrom
INSERT INTO schema_field_options (field_id, value, label, position)
VALUES ($1, $2, $3, $4);

-- name: GetFieldOptions :many
SELECT *
FROM schema_field_options
WHERE field_id = ANY($1::uuid[])
ORDER BY field_id, position;

-- name: DeleteFieldOptions :exec
DELETE FROM schema_field_options
WHERE field_id = $1::UUID;

-- name: CloneFieldOptions :execrows
INSERT INTO schema_field_options (field_id, value, label, position)
SELECT
    dst.object_id, -- new field's object_id
    src_opt.value,
    src_opt.label,
    src_opt.position
FROM schema_field_options src_opt
-- Get stable id from source field
         JOIN schema_fields src_f ON src_f.object_id = src_opt.field_id
-- Get new object_id from destination field with same stable id
         JOIN schema_fields dst ON dst.id = src_f.id
    AND dst.schema_version_id = sqlc.arg(draft_version_id)::UUID
WHERE src_f.schema_version_id = sqlc.arg(source_version_id)::UUID;

-- name: GetOptionByID :one
SELECT *
FROM schema_field_options
WHERE id = $1;

-- name: DeleteOptionByID :exec
DELETE FROM schema_field_options
WHERE id = $1;

-- name: CheckOptionValueInRules :one
SELECT EXISTS (
    SELECT 1 FROM schema_field_visibility_rules
    WHERE value @> jsonb_build_object('value', sqlc.arg(option_value)::text)
) OR EXISTS (
    SELECT 1 FROM schema_field_required_rules
    WHERE value @> jsonb_build_object('value', sqlc.arg(option_value)::text)
) AS is_referenced;

-- name: SetFieldOptions :exec
DELETE FROM schema_field_options
WHERE field_id = sqlc.arg(field_id)::UUID;

-- name: CreateFieldOptionForSet :exec
INSERT INTO schema_field_options (field_id, value, label, position)
VALUES ($1, $2, $3, $4);

-- /////////////////////////////// --
-- //// -- VISIBILITY RULES -- //// --
-- /////////////////////////////// --

-- name: CreateVisibilityRule :one
INSERT INTO schema_field_visibility_rules (
    field_id,
    depends_on_field_id,
    operator,
    value
)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: CreateVisibilityRulesBatch :copyfrom
INSERT INTO schema_field_visibility_rules (field_id, depends_on_field_id, operator, value)
VALUES ($1, $2, $3, $4);

-- name: GetFieldVisibilityRules :many
SELECT *
FROM schema_field_visibility_rules
WHERE field_id = ANY($1::uuid[])
ORDER BY field_id, created_at;

-- name: DeleteVisibilityRules :exec
DELETE FROM schema_field_visibility_rules
WHERE field_id = $1::UUID;

-- name: CloneVisibilityRules :execrows
INSERT INTO schema_field_visibility_rules (field_id, depends_on_field_id, operator, value)
SELECT
    dst_f.object_id,         -- new field's object_id
    dst_dep.object_id,       -- new dependent field's object_id
    src_rule.operator,
    src_rule.value
FROM schema_field_visibility_rules src_rule
-- Source field
         JOIN schema_fields src_f ON src_f.object_id = src_rule.field_id
-- Destination field (same stable id)
         JOIN schema_fields dst_f ON dst_f.id = src_f.id
    AND dst_f.schema_version_id = sqlc.arg(draft_version_id)::UUID
-- Source dependent field
JOIN schema_fields src_dep ON src_dep.object_id = src_rule.depends_on_field_id
-- Destination dependent field (same stable id)
    JOIN schema_fields dst_dep ON dst_dep.id = src_dep.id
    AND dst_dep.schema_version_id = sqlc.arg(draft_version_id)::UUID
WHERE src_f.schema_version_id = sqlc.arg(source_version_id)::UUID;

-- name: GetVisibilityRuleByID :one
SELECT *
FROM schema_field_visibility_rules
WHERE id = $1;

-- name: UpdateVisibilityRule :one
UPDATE schema_field_visibility_rules
SET
    depends_on_field_id = COALESCE(sqlc.arg(depends_on_field_id), depends_on_field_id),
    operator = COALESCE(sqlc.arg(operator), operator),
    value = COALESCE(sqlc.arg(value), value)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteVisibilityRuleByID :exec
DELETE FROM schema_field_visibility_rules
WHERE id = $1;

-- name: SetVisibilityRules :exec
DELETE FROM schema_field_visibility_rules
WHERE field_id = sqlc.arg(field_id)::UUID;

-- name: CreateVisibilityRuleForSet :one
INSERT INTO schema_field_visibility_rules (
    field_id,
    depends_on_field_id,
    operator,
    value
)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- /////////////////////////////// --
-- //// -- REQUIRED RULES -- //// --
-- /////////////////////////////// --

-- name: CreateRequiredRule :one
INSERT INTO schema_field_required_rules (
    field_id,
    depends_on_field_id,
    operator,
    value
)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: CreateRequiredRulesBatch :copyfrom
INSERT INTO schema_field_required_rules (field_id, depends_on_field_id, operator, value)
VALUES ($1, $2, $3, $4);

-- name: GetFieldRequiredRules :many
SELECT *
FROM schema_field_required_rules
WHERE field_id = ANY($1::uuid[])
ORDER BY field_id, created_at;

-- name: DeleteRequiredRules :exec
DELETE FROM schema_field_required_rules
WHERE field_id = $1::UUID;

-- name: CloneRequiredRules :execrows
INSERT INTO schema_field_required_rules (field_id, depends_on_field_id, operator, value)
SELECT
    dst_f.object_id,
    dst_dep.object_id,
    src_rule.operator,
    src_rule.value
FROM schema_field_required_rules src_rule
         JOIN schema_fields src_f ON src_f.object_id = src_rule.field_id
         JOIN schema_fields dst_f ON dst_f.id = src_f.id
    AND dst_f.schema_version_id = sqlc.arg(draft_version_id)::UUID
JOIN schema_fields src_dep ON src_dep.object_id = src_rule.depends_on_field_id
    JOIN schema_fields dst_dep ON dst_dep.id = src_dep.id
    AND dst_dep.schema_version_id = sqlc.arg(draft_version_id)::UUID
WHERE src_f.schema_version_id = sqlc.arg(source_version_id)::UUID;

-- name: GetRequiredRuleByID :one
SELECT *
FROM schema_field_required_rules
WHERE id = $1;

-- name: UpdateRequiredRule :one
UPDATE schema_field_required_rules
SET
    depends_on_field_id = COALESCE(sqlc.arg(depends_on_field_id), depends_on_field_id),
    operator = COALESCE(sqlc.arg(operator), operator),
    value = COALESCE(sqlc.arg(value), value)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRequiredRuleByID :exec
DELETE FROM schema_field_required_rules
WHERE id = $1;

-- name: SetRequiredRules :exec
DELETE FROM schema_field_required_rules
WHERE field_id = sqlc.arg(field_id)::UUID;

-- name: CreateRequiredRuleForSet :one
INSERT INTO schema_field_required_rules (
    field_id,
    depends_on_field_id,
    operator,
    value
)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- /////////////////////////////// --
-- //// -- FIELD CRUD -- //// --
-- /////////////////////////////// --

-- name: GetFieldByObjectID :one
SELECT *
FROM schema_fields
WHERE object_id = $1::UUID;

-- name: UpdateField :one
UPDATE schema_fields
SET
    key = COALESCE(sqlc.arg(key), key),
    type = COALESCE(sqlc.arg(type), type),
    title = COALESCE(sqlc.arg(title), title),
    description = COALESCE(sqlc.arg(description), description),
    placeholder = COALESCE(sqlc.arg(placeholder), placeholder),
    required = COALESCE(sqlc.arg(required), required),
    mutable = COALESCE(sqlc.arg(mutable), mutable),
    default_value = COALESCE(sqlc.arg(default_value), default_value),
    position = COALESCE(sqlc.arg(position), position),
    updated_at = NOW()
WHERE object_id = sqlc.arg(object_id)
  AND schema_version_id = sqlc.arg(schema_version_id)
RETURNING *;

-- name: CheckFieldKeyExists :one
SELECT EXISTS (
    SELECT 1
    FROM schema_fields
    WHERE schema_version_id = $1
      AND key = $2
      AND object_id != $3
) AS exists;

-- name: HasDependentRules :many
SELECT DISTINCT f.object_id as field_object_id, f.key as field_key
FROM schema_fields f
JOIN schema_field_visibility_rules vr ON vr.field_id = f.object_id
WHERE vr.depends_on_field_id = $1::UUID
   OR EXISTS (
       SELECT 1 FROM schema_field_required_rules rr
       WHERE rr.field_id = f.object_id
         AND rr.depends_on_field_id = $1::UUID
   );

-- name: DeleteField :exec
DELETE FROM schema_fields
WHERE object_id = $1::UUID;