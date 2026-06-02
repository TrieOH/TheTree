-- name: CreateField :one
INSERT INTO fields (step_id, key, title, description, position_hint, required, type, placeholder, default_value, config)
VALUES (@step_id, @key, @title, @description, @position_hint, @required, @type, @placeholder, @default_value, @config)
RETURNING *;

-- name: GetFieldByID :one
SELECT * FROM fields
WHERE id = $1;

-- name: ListFieldsByStepID :many
SELECT * FROM fields
WHERE step_id = $1
ORDER BY position_hint ASC;

-- name: ListFieldsByFormID :many
SELECT f.* FROM fields f
INNER JOIN steps s ON f.step_id = s.id
WHERE s.form_id = @form_id
ORDER BY f.position_hint ASC;

-- name: BulkEditFields :batchexec
UPDATE fields
SET key           = @key,
    title         = @title,
    description   = @description,
    position_hint = @position_hint,
    required      = @required,
    type          = @type,
    placeholder   = @placeholder,
    default_value = @default_value,
    config        = @config,
    updated_at    = now()
WHERE id = @id
  AND step_id = @step_id;

-- name: DeleteField :exec
DELETE FROM fields
WHERE id = $1;

-- name: CreateFieldSelectConfig :one
INSERT INTO field_select_config (field_id, behaviour, value_type, options)
VALUES (@field_id, @behaviour, @value_type, @options)
RETURNING *;

-- name: GetFieldSelectConfig :one
SELECT * FROM field_select_config
WHERE field_id = $1;

-- name: UpdateFieldSelectConfig :one
UPDATE field_select_config
SET behaviour  = @behaviour,
    value_type = @value_type,
    options    = @options
WHERE field_id = @field_id
RETURNING *;

-- name: DeleteFieldSelectConfig :exec
DELETE FROM field_select_config
WHERE field_id = $1;