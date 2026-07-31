-- +goose Up
CREATE TABLE products (
    id                    UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id            UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    vendor_code           VARCHAR(256) NOT NULL,
    requires_registration BOOLEAN NOT NULL DEFAULT TRUE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ,
    deleted_at            TIMESTAMPTZ,
    CONSTRAINT uniq_products_edition_vendor_code UNIQUE(edition_id, vendor_code)
);

CREATE TABLE product_variants (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id   UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    product_id   UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    vendor_code  VARCHAR(256) NOT NULL,
    name         VARCHAR(128) NOT NULL,
    description  TEXT,
    price        BIGINT NOT NULL DEFAULT 0, -- cents
    stock        INT, -- null = unlimited
    gallery_urls JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ,
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT uniq_product_variants_edition_vendor_code UNIQUE(edition_id, vendor_code)
);

CREATE TABLE product_purchases (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id        UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    variant_id        UUID NOT NULL REFERENCES product_variants(id),
    purchaser_id      UUID NOT NULL,
    recipient_id      UUID,
    registration_id   UUID REFERENCES registrations(id),
    quantity          INT NOT NULL DEFAULT 1,
    status            TEXT NOT NULL DEFAULT 'pending',
    CONSTRAINT chk_product_purchases_status_valid CHECK (
        status IN ('pending', 'confirmed', 'cancelled', 'expired')
    ),
    status_reason     TEXT,
    payssage_intent_id UUID,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);
-- +goose Down
DROP TABLE IF EXISTS product_purchases;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;