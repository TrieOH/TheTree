-- +goose Up
-- The profile-badges gift claim scans registrations by attendee_email across
-- all editions (claim-all-by-email), unlike the edition-scoped my-ticket
-- claim which the per-edition partial index (013) covers.
CREATE INDEX idx_registrations_attendee_email
    ON registrations (attendee_email);
-- +goose Down
DROP INDEX IF EXISTS idx_registrations_attendee_email;
