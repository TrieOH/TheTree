-- name: GetRegistrationByID :one
SELECT *
FROM registrations
WHERE id = @id;

-- name: CreateRegistration :one
INSERT INTO registrations (edition_id, ticket_type_id, purchaser_id, attendee_user_id, attendee_email, attendee_name, status, status_reason, payssage_intent_id)
VALUES (@edition_id, @ticket_type_id, @purchaser_id, @attendee_user_id, @attendee_email, @attendee_name, @status, @status_reason, @payssage_intent_id)
RETURNING *;

-- name: UpdateRegistrationStatus :one
UPDATE registrations
SET
    status        = @status,
    status_reason = @status_reason,
    updated_at    = now()
WHERE id = @id
RETURNING *;
