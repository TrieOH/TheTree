-- name: UpsertBadgeEmission :one
INSERT INTO badge_emissions (edition_id, user_id, origin, registration_id)
VALUES (@edition_id, @user_id, @origin, @registration_id)
ON CONFLICT (edition_id, user_id, origin) DO UPDATE
SET
    status        = 'active',
    status_reason = NULL,
    updated_at    = now()
RETURNING *;

-- name: RevokeBadgeEmission :exec
UPDATE badge_emissions
SET
    status        = 'revoked',
    status_reason = @reason,
    updated_at    = now()
WHERE edition_id = @edition_id
  AND user_id = @user_id
  AND origin = @origin
  AND status = 'active';

-- name: RevokeBadgeEmissionByRegistration :exec
UPDATE badge_emissions
SET
    status        = 'revoked',
    status_reason = @reason,
    updated_at    = now()
WHERE registration_id = @registration_id
  AND origin = 'participant'
  AND status = 'active';

-- name: MarkBadgeEmissionEmailSent :exec
UPDATE badge_emissions
SET
    email_sent_at = now(),
    updated_at    = now()
WHERE id = @id;

-- name: ListBadgeEmissionViewsByUser :many
SELECT
    be.id,
    be.edition_id,
    be.user_id,
    be.origin,
    be.registration_id,
    be.status,
    be.status_reason,
    be.email_sent_at,
    be.emitted_at,
    be.updated_at,
    e.edition_name,
    e.ends_at,
    ev.full_name AS event_name,
    r.ticket_type_id,
    tt.name AS ticket_name
FROM badge_emissions be
JOIN editions e ON e.id = be.edition_id AND e.deleted_at IS NULL
JOIN events ev ON ev.id = e.event_id AND ev.deleted_at IS NULL
LEFT JOIN registrations r ON r.id = be.registration_id AND r.deleted_at IS NULL
LEFT JOIN ticket_types tt ON tt.id = r.ticket_type_id AND tt.deleted_at IS NULL
WHERE be.user_id = @user_id
  AND be.status = 'active'
ORDER BY e.ends_at DESC, be.emitted_at DESC;

-- name: ListBadgeEmissionViewsByEdition :many
SELECT
    be.id,
    be.edition_id,
    be.user_id,
    be.origin,
    be.registration_id,
    be.status,
    be.status_reason,
    be.email_sent_at,
    be.emitted_at,
    be.updated_at,
    e.edition_name,
    e.ends_at,
    ev.full_name AS event_name,
    r.ticket_type_id,
    tt.name AS ticket_name
FROM badge_emissions be
JOIN editions e ON e.id = be.edition_id AND e.deleted_at IS NULL
JOIN events ev ON ev.id = e.event_id AND ev.deleted_at IS NULL
LEFT JOIN registrations r ON r.id = be.registration_id AND r.deleted_at IS NULL
LEFT JOIN ticket_types tt ON tt.id = r.ticket_type_id AND tt.deleted_at IS NULL
WHERE be.edition_id = @edition_id
ORDER BY be.user_id, be.origin;
