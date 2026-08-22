-- +goose Up
-- One active participation per attendee per occurrence. Paid checkouts and
-- self-service registrations both insert; this partial unique index is the
-- concurrency backstop — two racing registrations for the same
-- (occurrence, registration) can only lose here, the second insert violates
-- the index atomically (mirrors uniq_registrations_active_per_edition_attendee).
-- Cancelled rows free the slot; re-registering after a cancel inserts a fresh
-- row (append-only ledger, like ticket registrations).
CREATE UNIQUE INDEX uniq_program_participations_active_per_occurrence_attendee
    ON program_participations (occurrence_id, registration_id)
    WHERE status IN ('registered', 'attended', 'no_show');

-- +goose Down
DROP INDEX IF EXISTS uniq_program_participations_active_per_occurrence_attendee;
