-- 012_one_ticket_per_attendee.sql
-- +goose Up
-- One active ticket per person per edition. An attendee holds a ticket
-- through their registration (pending = reserved, confirmed = paid);
-- cancelled/expired registrations free the slot. This is the concurrency
-- backstop: the checkout service pre-checks (clean 409) inside the tx, but
-- two racing checkouts both inserting for the same attendee can only lose
-- here — the second insert violates the index atomically.
CREATE UNIQUE INDEX uniq_registrations_active_per_edition_attendee
    ON registrations (edition_id, attendee_user_id)
    WHERE status IN ('pending', 'confirmed') AND deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS uniq_registrations_active_per_edition_attendee;
