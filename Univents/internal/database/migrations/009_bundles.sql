-- +goose Up

-- Product bundles (composite products)
CREATE TABLE product_bundle_items (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE, -- the bundle parent

    item_type VARCHAR(32) NOT NULL CHECK (item_type IN ('ticket', 'product', 'tokens')),
    item_id UUID NULL, -- ticket_id or product_id, null for tokens
    token_amount INT NOT NULL DEFAULT 0, -- if item_type = 'tokens'

    quantity INT NOT NULL DEFAULT 1,
    included_price_cents INT NOT NULL DEFAULT 0, -- override price (0 = use item's price)

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(product_id, item_type, item_id, token_amount)
);

CREATE INDEX idx_bundle_items_product ON product_bundle_items(product_id);

-- Token products (forced base product)
CREATE TABLE edition_token_settings (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL UNIQUE REFERENCES editions(id) ON DELETE CASCADE,

    token_product_id UUID NULL REFERENCES products(id), -- auto-created, editable

    -- organizer customization
    token_price_cents INT NOT NULL DEFAULT 0, -- price per 1 token
    min_purchase_amount INT NOT NULL DEFAULT 1,
    max_purchase_amount INT NOT NULL DEFAULT 100,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger: when edition.has_tokens becomes true, create token product if missing
-- Trigger: when token_price_cents changes, update product price

-- +goose Down
DROP TABLE IF EXISTS edition_token_settings;
DROP INDEX IF EXISTS idx_bundle_items_product;
DROP TABLE IF EXISTS product_bundle_items;