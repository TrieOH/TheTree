-- +goose Up
CREATE TABLE ticket_types (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id   UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    name         VARCHAR(128) NOT NULL,
    description  TEXT,
    access_level INT NOT NULL DEFAULT 0,
    -- cents, 0 = free
    price        INT NOT NULL DEFAULT 0,
    -- null = unlimited
    max_quantity INT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ
);
-- +goose Down
DROP TABLE IF EXISTS ticket_types;