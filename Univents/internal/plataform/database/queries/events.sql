-- name: CreateEvent :one
INSERT INTO events (id, owner_id, organization_id, name, acronym, slug, tagline, description, is_series, logo_url, contact_email, social_links, created_by, goauth_scope_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: PatchEvent :one
UPDATE events
SET
    name = $1,
    acronym = $2,
    slug = $3,
    tagline = $4,
    description = $5,
    is_series = $6,
    logo_url = $7,
    banner_url = $8,
    has_gallery = $9,
    contact_email = $10,
    social_links = $11,
    updated_at = now(),
    owner_id = $14
WHERE
    id = $12
    AND goauth_scope_id = $13
    AND deleted_at IS NULL
    RETURNING *;

-- name: GetEventByID :one
SELECT *
FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT *
FROM events
WHERE status != 'draft';

-- name: ListOwnEvents :many
SELECT *
FROM events
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: PublishEvent :exec
UPDATE events
SET status = 'active'
WHERE id = $1;

-- name: AddEdition :execrows
UPDATE events
SET editions_count = editions_count + 1
WHERE id = $1
    AND (
        is_series = true
        OR editions_count = 0
    );

-- name: AddEventGalleryImage :one
UPDATE events
SET gallery_urls = array_append(COALESCE(gallery_urls, '{}'), @url::text)
WHERE id = @id
  AND deleted_at IS NULL
    RETURNING *;

-- name: RemoveEventGalleryImage :one
UPDATE events
SET gallery_urls = array_remove(gallery_urls, @url::text)
WHERE id = @id
  AND deleted_at IS NULL
    RETURNING *;

-- name: SetEventLogo :one
UPDATE events
SET
    logo_url = @url::text
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: UnsetEventLogo :one
UPDATE events
SET logo_url = NULL
WHERE id = @id
  AND deleted_at IS NULL
    RETURNING *;

-- name: SetEventBanner :one
UPDATE events
SET
    banner_url = @url::text
WHERE id = @id
  AND deleted_at IS NULL
    RETURNING *;

-- name: UnsetEventBanner :one
UPDATE events
SET banner_url = NULL
WHERE id = @id
  AND deleted_at IS NULL
    RETURNING *;

-------------------------------------------------------------------------------------------
------------------------------------------ Audit ------------------------------------------
-------------------------------------------------------------------------------------------

-- name: AppendEventAudits :copyfrom
INSERT INTO event_audit (event_id, actor_type, actor_id, action, state, from_status, to_status, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListEventAuditByEvent :many
SELECT * FROM event_audit
WHERE event_id = $1
ORDER BY created_at DESC;

-- name: ListEventAuditByEventAndAction :many
SELECT * FROM event_audit
WHERE event_id = $1 AND action = $2
ORDER BY created_at DESC
    LIMIT $3 OFFSET $4;

-- name: GetEventAudit :one
SELECT * FROM event_audit WHERE id = $1;

-- name: GetLastEventStatusChange :one
SELECT * FROM event_audit
WHERE event_id = $1
  AND from_status IS NOT NULL
  AND to_status IS NOT NULL
ORDER BY created_at DESC
    LIMIT 1;