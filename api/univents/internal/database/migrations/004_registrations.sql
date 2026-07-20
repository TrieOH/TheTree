-- +goose Up
CREATE TABLE registrations (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id       UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    ticket_type_id   UUID NOT NULL REFERENCES ticket_types(id),
    purchaser_id     UUID NOT NULL,
    attendee_user_id UUID,
    attendee_email   VARCHAR(256) NOT NULL,
    attendee_name    VARCHAR(256),
    status           TEXT NOT NULL DEFAULT 'pending',
    CONSTRAINT chk_registrations_status_valid CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'expired')
    ),
    status_reason    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ
);
-- +goose Down
DROP TABLE IF EXISTS registrations;