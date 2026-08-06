-- name: CreateBadgeTemplate :one
INSERT INTO badge_templates (
    edition_id, ticket_type_id, origin, name, design_data
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBadgeTemplateByID :one
SELECT * FROM badge_templates
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListBadgeTemplatesByEdition :many
SELECT * FROM badge_templates
WHERE edition_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteBadgeTemplate :exec
UPDATE badge_templates
SET deleted_at = now(),
    updated_at = now()
WHERE id = $1;

-- name: UpdateBadgeTemplate :one
UPDATE badge_templates
SET
    name        = @name,
    design_data = @design_data,
    updated_at  = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: ListBadgeTemplatesByEditionIDs :many
SELECT *
FROM badge_templates
WHERE edition_id = ANY(@edition_ids::uuid[])
  AND deleted_at IS NULL;
