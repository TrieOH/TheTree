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

-- name: GetEditionByID :one
SELECT *
FROM editions
WHERE id = @id;

-- name: GetEditionBySlug :one
SELECT *
FROM editions
WHERE slug = @slug
  AND event_id = @event_id
  AND is_draft = FALSE;

-- name: PublishEdition :exec
UPDATE editions
SET
    is_draft = FALSE,
    updated_at = now()
WHERE id = @id
  AND is_draft = TRUE;

-- name: PatchEdition :one
UPDATE editions
SET
    edition_name           = @edition_name,
    slug                   = @slug,
    tagline                = @tagline,
    description            = @description,
    registration_opens_at  = @registration_opens_at,
    starts_at              = @starts_at,
    ends_at                = @ends_at,
    location_name          = @location_name,
    location_address       = @location_address,
    logo_url               = @logo_url,
    banner_url             = @banner_url,
    contact_email          = @contact_email,
    updated_at             = now()
WHERE id = @id
  AND is_draft = TRUE
RETURNING *;