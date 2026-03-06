-- +goose Up
-- Products (merchandise, addons, token packs)
CREATE TYPE product_type AS ENUM (
    'merchandise',  -- shirts, stickers, etc
    'ticket',       -- event access
    'token',        -- bundles of tokens
    'bundle'        -- composite product
);

CREATE TYPE product_status AS ENUM (
    'draft',
    'available',
    'sold_out',
    'unavailable'
);

-- CREATE TYPE bundle_discount_type AS ENUM (
--    'percentage',
--    'fixed_amount',
--    'none'
--);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    -- identity
    name VARCHAR(256) NOT NULL,
    description TEXT NULL,
    type product_type NOT NULL,
    ticket_id UUID NULL REFERENCES tickets(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CHECK (type != 'ticket' OR ticket_id IS NOT NULL),

    -- pricing
    price_cents INT NOT NULL DEFAULT 0, -- 0 = free

    -- visibility
    status product_status NOT NULL DEFAULT 'draft',
    available_from TIMESTAMPTZ NULL,
    available_until TIMESTAMPTZ NULL,

    -- inventory
    has_inventory BOOLEAN NOT NULL DEFAULT FALSE,
    inventory_quantity INT NOT NULL DEFAULT 0,
    inventory_remaining INT NOT NULL DEFAULT 0,

    CHECK (has_inventory = FALSE OR inventory_quantity > 0),
    CHECK (inventory_quantity >= inventory_remaining),

    -- bundles
    -- bundle_display_comparison BOOLEAN NOT NULL DEFAULT TRUE, -- show "was $X, now $Y"
    -- bundle_discount_type bundle_discount_type NULL,
    -- bundle_discount_value INT NULL,

    -- purchase limits
    -- has_user_limit BOOLEAN NOT NULL DEFAULT FALSE,
    -- max_per_user INT NOT NULL DEFAULT 1,
    -- min_per_order INT NOT NULL DEFAULT 1,

    -- tokens (if type = token)
    -- has_tokens_included BOOLEAN NOT NULL DEFAULT FALSE,
    -- tokens_included INT NOT NULL DEFAULT 0,
    -- CONSTRAINT chk_tokens_if_pack
    --    CHECK (type != 'token' OR tokens_included > 0),

    -- media
    -- has_thumbnail BOOLEAN NOT NULL DEFAULT FALSE,
    -- thumbnail_url TEXT NULL,
    -- has_gallery BOOLEAN NOT NULL DEFAULT FALSE,
    -- gallery_image_urls TEXT[] NULL,

    -- fulfillment methods (at least one must be true)
    -- has_event_pickup BOOLEAN NOT NULL DEFAULT FALSE,      -- collect at venue
    -- has_physical_shipping BOOLEAN NOT NULL DEFAULT FALSE, -- mail delivery
    -- has_digital_delivery BOOLEAN NOT NULL DEFAULT FALSE,  -- instant download
    -- has_auto_fulfillment BOOLEAN NOT NULL DEFAULT FALSE,  -- tokens, tier upgrades, virtual goods

    -- CONSTRAINT chk_one_fulfillment_method
    --    CHECK (has_event_pickup OR has_physical_shipping OR has_digital_delivery OR has_auto_fulfillment),

    -- shipping config
    -- shipping_required_at_checkout BOOLEAN NOT NULL DEFAULT FALSE,
    -- shipping_address_fields JSONB NULL, -- {country: true, phone: true}

    -- digital config
    -- digital_content JSONB NULL, -- {download_url: "...", expires_at: "...", access_code: "..."}

    -- pickup config
    -- pickup_location VARCHAR(256) NULL, -- "Registration Desk", "VIP Lounge"
    -- pickup_window_start TIMESTAMPTZ NULL,
    -- pickup_window_end TIMESTAMPTZ NULL,

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_products_edition ON products(edition_id, status, type)
    WHERE deleted_at IS NULL;

-- Bundle items (what's inside)
CREATE TABLE product_bundle_components (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    bundle_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    component_type VARCHAR(32) NOT NULL CHECK (component_type IN ('ticket', 'product')),
    component_id UUID NOT NULL, -- ticket_id or product_id

    quantity INT NOT NULL DEFAULT 1,

    -- pricing override (optional)
    use_component_price BOOLEAN NOT NULL DEFAULT TRUE, -- false = custom price for this bundle
    override_price_cents INT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(bundle_product_id, component_type, component_id)
);

CREATE TABLE product_reservations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    session_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(session_id, product_id)
);

CREATE INDEX idx_reservations_session ON product_reservations(session_id);
CREATE INDEX idx_reservations_expires ON product_reservations(expires_at); -- for cleanup jobs

-- Purchases (unified for tickets + products)
CREATE TYPE purchase_status AS ENUM (
    'pending',        -- cart reserved, unpaid
    'completed',      -- paid, fulfilled
    'refunded',       -- fully refunded
    'partial_refund',
    'cancelled'
);

CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,

    status purchase_status NOT NULL DEFAULT 'pending',

    -- totals
    subtotal_cents INT NOT NULL DEFAULT 0,
    discount_cents INT NOT NULL DEFAULT 0,
    tax_cents INT NOT NULL DEFAULT 0,
    tax_breakdown JSONB NULL,
    total_cents INT NOT NULL DEFAULT 0,

    -- payment ref (integrate with your payment provider)
    payment_provider VARCHAR(32) NULL,
    payment_id TEXT NULL,

    -- fulfillment
    fulfilled_at TIMESTAMPTZ NULL,
    fulfilment_notes TEXT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_purchases_user ON purchases(user_id, status);
CREATE INDEX idx_purchases_edition ON purchases(edition_id, created_at);

-- Purchase items (line items)
CREATE TABLE purchase_items (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,

    item_type VARCHAR(32) NOT NULL, -- 'ticket', 'product'
    item_id UUID NOT NULL, -- ticket_id or product_id

    quantity INT NOT NULL DEFAULT 1,
    unit_price_cents INT NOT NULL,
    total_price_cents INT NOT NULL,

    -- gifting
    assigned_to_user_id UUID NULL, -- if buying for someone else

    -- product-specific
    fulfilled BOOLEAN NOT NULL DEFAULT FALSE,
    fulfilled_at TIMESTAMPTZ NULL,

    -- refund tracking
    refunded_quantity INT NOT NULL DEFAULT 0,
    refunded_amount_cents INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_purchase_items_purchase ON purchase_items(purchase_id);
CREATE INDEX idx_purchase_items_assigned ON purchase_items(assigned_to_user_id)
    WHERE assigned_to_user_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_purchase_items_assigned;
DROP INDEX IF EXISTS idx_purchase_items_purchase;
DROP TABLE IF EXISTS purchase_items;
DROP INDEX IF EXISTS idx_purchases_edition;
DROP INDEX IF EXISTS idx_purchases_user;
DROP TABLE IF EXISTS purchases;
DROP TYPE IF EXISTS purchase_status;
DROP INDEX IF EXISTS idx_products_edition;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS product_status;
DROP TYPE IF EXISTS product_type;