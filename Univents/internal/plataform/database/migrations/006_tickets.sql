-- +goose Up
-- Tickets (access control)
CREATE TYPE ticket_status AS ENUM (
    'draft',
    'on_sale',
    'paused',
    'sold_out',
    'unavailable'
);

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    -- identity
    name VARCHAR(256) NOT NULL,
    description TEXT NULL,

    -- pricing
    price_cents INT NOT NULL DEFAULT 0, -- 0 = free ticket
    
    has_limited_quantity BOOLEAN NOT NULL DEFAULT FALSE,
    quantity_available INT NOT NULL,

    quantity_sold INT NOT NULL DEFAULT 0,
    quantity_reserved INT NOT NULL DEFAULT 0, -- pending checkouts

    CHECK (has_limited_quantity OR (quantity_available = 0 AND quantity_sold = 0)),
    CHECK (quantity_sold + quantity_reserved <= quantity_available OR NOT has_limited_quantity),

    -- visibility
    status ticket_status NOT NULL DEFAULT 'draft',

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_tickets_edition ON tickets(edition_id, status)
    WHERE deleted_at IS NULL;

CREATE TYPE permission_type AS ENUM (
    'activity',
    'product',
    'checkpoint'
);

-- Ticket permissions (what each tier can access)
CREATE TABLE ticket_permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,

    permission_type permission_type NOT NULL, -- 'activity', 'product', 'checkpoint', 'zone'
    target_id UUID NULL, -- specific activity/product ID, null = all of type

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(ticket_id, permission_type, target_id)
);

-- User ticket grants (effective access after purchase)
CREATE TABLE user_tickets (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,

    -- source
    purchase_item_id UUID NOT NULL REFERENCES purchase_items(id),
    ticket_id UUID NOT NULL REFERENCES tickets(id),

    -- status
    status VARCHAR(32) NOT NULL DEFAULT 'active', -- 'active', 'used', 'transferred', 'refunded'

    -- usage
    checked_in BOOLEAN NOT NULL DEFAULT FALSE,
    checked_in_at TIMESTAMPTZ NULL,

    -- transfer
    transferred_to_user_id UUID NULL,
    transferred_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(purchase_item_id) -- one grant per purchase item
);

CREATE INDEX idx_user_tickets_user ON user_tickets(user_id, edition_id, status);
CREATE INDEX idx_user_tickets_ticket ON user_tickets(ticket_id, status);

-- +goose Down
DROP INDEX IF EXISTS idx_user_tickets_ticket;
DROP INDEX IF EXISTS idx_user_tickets_user;
DROP TABLE IF EXISTS user_tickets;
DROP TABLE IF EXISTS ticket_tier_permissions;
DROP INDEX IF EXISTS idx_tickets_edition;
DROP TABLE IF EXISTS tickets;
DROP TYPE IF EXISTS ticket_status;

