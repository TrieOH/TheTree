-- name: CreateEvent :one
INSERT INTO events (owner_id, full_name, acronym, slug, description, status, contact_email)
VALUES (@owner_id, @full_name, @acronym, @slug, @description, @status, @contact_email)
RETURNING *;

-- name: GetEventByID :one
SELECT *
FROM events
WHERE id = @id;

-- name: ListPublicEvents :many
SELECT *
FROM events
WHERE status != 'draft';

-- name: GetEventBySlug :one
SELECT *
FROM events
WHERE slug = @slug
  AND status != 'draft';

-- name: ListOwnedEvents :many
SELECT *
FROM events
WHERE owner_id = @owner_id;

-- name: ListJoinedEvents :many
SELECT e.*
FROM events e
JOIN event_members em ON em.event_id = e.id
WHERE em.user_id = @user_id
  AND em.deleted_at IS NULL
  AND e.deleted_at IS NULL;

-- name: PublishEvent :exec
UPDATE events
SET
    status = 'active',
    updated_at = now()
WHERE id = @id
  AND status = 'draft';

-- name: DiscontinueEvent :exec
UPDATE events
SET
    status = 'discontinued',
    updated_at = now()
WHERE id = @id
  AND status = 'active';

-- name: PatchEvent :one
UPDATE events
SET
    full_name     = @full_name,
    acronym       = @acronym,
    slug          = @slug,
    description   = @description,
    logo_url      = @logo_url,
    banner_url    = @banner_url,
    contact_email = @contact_email,
    updated_at    = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: GetEventMember :one
SELECT *
FROM event_members
WHERE event_id = @event_id
  AND user_id = @user_id
  AND deleted_at IS NULL;

-- name: AddEventMember :one
INSERT INTO event_members (event_id, user_id, role)
VALUES (@event_id, @user_id, @role)
RETURNING *;

-- name: RemoveEventMember :exec
UPDATE event_members
SET
    deleted_at = now(),
    updated_at = now()
WHERE event_id = @event_id
  AND user_id = @user_id
  AND deleted_at IS NULL;

-- name: ListEventMembers :many
SELECT *
FROM event_members
WHERE event_id = @event_id
  AND deleted_at IS NULL
ORDER BY created_at;