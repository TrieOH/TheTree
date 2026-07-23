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