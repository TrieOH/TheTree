-- +goose Up
-- Tickets (access control)
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    -- identity
    name VARCHAR(256) NOT NULL,
    description TEXT NULL,

    UNIQUE (name, edition_id),

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

-- User ticket grants (effective access after purchase)
--CREATE TABLE user_tickets (
--    id UUID PRIMARY KEY DEFAULT uuidv7(),
--    user_id UUID NOT NULL,
--    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,

    -- source
    --purchase_item_id UUID NOT NULL REFERENCES purchase_items(id),
    --ticket_id UUID NOT NULL REFERENCES tickets(id),

    -- status
    --status VARCHAR(32) NOT NULL DEFAULT 'active', -- 'active', 'used', 'transferred', 'refunded'

    -- usage
    --checked_in BOOLEAN NOT NULL DEFAULT FALSE,
    --checked_in_at TIMESTAMPTZ NULL,

    -- transfer
    --transferred_to_user_id UUID NULL,
    --transferred_at TIMESTAMPTZ NULL,

    --created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    --updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    --UNIQUE(purchase_item_id) -- one grant per purchase item
--);

--CREATE INDEX idx_user_tickets_user ON user_tickets(user_id, edition_id, status);
--CREATE INDEX idx_user_tickets_ticket ON user_tickets(ticket_id, status);

-- +goose Down
--DROP INDEX IF EXISTS idx_user_tickets_ticket;
--DROP INDEX IF EXISTS idx_user_tickets_user;
--DROP TABLE IF EXISTS user_tickets;
DROP TABLE IF EXISTS tickets;

