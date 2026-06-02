-- name: CreateResponse :one
INSERT INTO responses (
    form_id,
    invite_id,
    responder_id,
    email
) VALUES (
    @form_id,
    @invite_id,
    @responder_id,
    @email
)
RETURNING *;

-- name: FinishResponse :exec
UPDATE responses
SET finished_at = NOW()
WHERE id = @id
  AND finished_at IS NULL;

-- name: GetResponseByID :one
SELECT * FROM responses
WHERE id = @id;

-- name: ListResponsesByForm :many
SELECT * FROM responses
WHERE form_id = @form_id
  AND finished_at IS NOT NULL
ORDER BY started_at DESC;