-- name: GetRegistrationByID :one
SELECT *
FROM registrations
WHERE id = @id;

-- name: CreateRegistration :one
INSERT INTO registrations (edition_id, ticket_type_id, purchaser_id, attendee_user_id, attendee_email, attendee_name, status, status_reason, payssage_intent_id)
VALUES (@edition_id, @ticket_type_id, @purchaser_id, @attendee_user_id, @attendee_email, @attendee_name, @status, @status_reason, @payssage_intent_id)
RETURNING *;

-- name: GetActiveByEditionAndAttendee :one
-- The attendee's active (pending or confirmed) registration in an edition,
-- if any — the one-ticket-per-person check (checkout) and the my-ticket
-- read. Cancelled/expired registrations do not count.
SELECT *
FROM registrations
WHERE edition_id = @edition_id
  AND attendee_user_id = @attendee_user_id
  AND status IN ('pending', 'confirmed')
  AND deleted_at IS NULL
LIMIT 1;

-- name: CountConfirmedByEdition :one
-- The number of confirmed attendees of an edition — the attendee-count
-- read (GET /editions/{id}/attendees/count). Only paid registrations
-- count; pending reservations and cancelled/expired rows do not.
SELECT COUNT(*)
FROM registrations
WHERE edition_id = @edition_id
  AND status = 'confirmed'
  AND deleted_at IS NULL;

-- name: UpdateRegistrationStatus :one
UPDATE registrations
SET
    status        = @status,
    status_reason = @status_reason,
    updated_at    = now()
WHERE id = @id
RETURNING *;
