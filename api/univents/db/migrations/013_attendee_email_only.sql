-- +goose Up
-- Gifted tickets may target a recipient with no IdentityX account yet: the
-- checkout resolves the recipient email against IdentityX and ties the actor
-- id when the account exists; otherwise the registration is created
-- email-only (attendee_user_id NULL) and the recipient claims it after
-- creating an account under that email. Every existing row has a user id, so
-- this only relaxes the constraint for new email-only gifts.
ALTER TABLE registrations ALTER COLUMN attendee_user_id DROP NOT NULL;

-- One active ticket per person per edition, by email (PG18 NULLS NOT
-- DISTINCT): email is the always-present attendee identity, so the same
-- index covers accountless recipients (NULL attendee_user_id) and account
-- holders alike — two racing checkouts gifting to the same person can only
-- lose here, the second insert violates the index atomically. (Naively
-- adding NULLS NOT DISTINCT to the attendee_user_id index would be wrong:
-- two *different* accountless recipients both have NULL and would collide.)
-- The attendee_user_id index stays as-is: it remains the dedup for account
-- holders whose email changed between purchases.
CREATE UNIQUE INDEX uniq_registrations_active_email_per_edition
    ON registrations (edition_id, attendee_email) NULLS NOT DISTINCT
    WHERE status IN ('pending', 'confirmed') AND deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS uniq_registrations_active_email_per_edition;
ALTER TABLE registrations ALTER COLUMN attendee_user_id SET NOT NULL;
