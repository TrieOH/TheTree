-- name: CreateEdition :one
INSERT INTO editions (event_id, edition_name, slug, starts_at, ends_at, created_by)
VALUES (@event_id, @edition_name, @slug, @starts_at, @ends_at, @created_by)
RETURNING *;

-- name: ListPublicEditions :many
SELECT *
FROM editions
WHERE is_draft = FALSE
  AND event_id = @event_id;

-- name: ListDraftEditions :many
SELECT *
FROM editions
WHERE is_draft = TRUE
  AND event_id = @event_id;