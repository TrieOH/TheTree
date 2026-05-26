-- name: CreateStep :one
INSERT INTO steps (form_id, title, description, position_hint)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: ListStepsByFormID :many
SELECT * FROM steps
WHERE form_id = $1
ORDER BY position_hint ASC;

-- name: BulkEditSteps :batchexec
UPDATE steps
SET title = @title,
    description = @description,
    position_hint = @position_hint
WHERE id = @id
  AND form_id = @form_id;