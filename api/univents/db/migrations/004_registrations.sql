-- +goose Up
CREATE TABLE registrations (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id       UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    ticket_type_id   UUID NOT NULL REFERENCES ticket_types(id),
    purchaser_id     UUID NOT NULL,
    attendee_user_id UUID NOT NULL,
    attendee_email   VARCHAR(256) NOT NULL,
    attendee_name    VARCHAR(256) NOT NULL ,
    status           TEXT NOT NULL DEFAULT 'pending',
    CONSTRAINT chk_registrations_status_valid CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'expired')
    ),
    status_reason    TEXT,
    payssage_intent_id UUID,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ
);

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
DROP TABLE IF EXISTS registrations;