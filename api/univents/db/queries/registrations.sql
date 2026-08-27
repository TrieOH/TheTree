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

-- name: GetActiveByEditionAndAttendeeEmail :one
-- The email-only recipient's active (pending or confirmed) registration in
-- an edition — the one-ticket-per-person pre-check for gifted tickets to
-- recipients without an IdentityX account (attendee_user_id NULL). Only
-- rows with a NULL user id match: a recipient with an account is resolved
-- to a user id at checkout and covered by GetActiveByEditionAndAttendee.
-- Cancelled/expired registrations do not count.
SELECT *
FROM registrations
WHERE edition_id = @edition_id
  AND attendee_email = @attendee_email
  AND attendee_user_id IS NULL
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

-- name: ClaimRegistrationByEmail :one
-- Ties an email-only gifted registration (attendee_user_id NULL) to the
-- recipient's newly created IdentityX account — the gift claim, fired
-- lazily from the my-ticket read when the caller holds no ticket under
-- their id. Only active rows match (pending = paid-but-unconfirmed,
-- confirmed = badge already deferred); the active-per-edition email index
-- guarantees at most one match. Cancelled/expired gifts are history.
UPDATE registrations
SET
    attendee_user_id = @attendee_user_id,
    updated_at       = now()
WHERE edition_id = @edition_id
  AND attendee_email = @attendee_email
  AND attendee_user_id IS NULL
  AND status IN ('pending', 'confirmed')
  AND deleted_at IS NULL
RETURNING *;

-- name: ClaimAllByAttendeeEmail :many
-- Ties every active email-only gift for this email to the recipient's
-- IdentityX account — the profile-badges claim, fired from the badges
-- read (a profile requires an account, so the account email
-- deterministically resolves the gifts). Crosses editions, unlike
-- ClaimRegistrationByEmail. Returns the claimed registrations so the
-- caller can emit deferred badges for the confirmed ones; pending gifts
-- are claimed too (tied to the account as soon as it exists).
-- Cancelled/expired gifts are history.
UPDATE registrations
SET
    attendee_user_id = @attendee_user_id,
    updated_at       = now()
WHERE attendee_email = @attendee_email
  AND attendee_user_id IS NULL
  AND status IN ('pending', 'confirmed')
  AND deleted_at IS NULL
RETURNING *;
