-- +goose Up
ALTER TABLE products
    ADD COLUMN thumbnail_url TEXT NULL,
    ADD COLUMN gallery_urls TEXT[] NULL,
    ADD COLUMN hard_deleted_at TIMESTAMPTZ NULL;

ALTER TYPE product_status ADD VALUE 'hidden';

CREATE INDEX idx_products_pending_hard_delete
    ON products(deleted_at)
    WHERE deleted_at IS NOT NULL AND hard_deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_products_pending_hard_delete;
ALTER TABLE products
DROP COLUMN IF EXISTS thumbnail_url,
    DROP COLUMN IF EXISTS gallery_urls,
    DROP COLUMN IF EXISTS hard_deleted_at;

-- NOTE: 'hidden' cannot be removed from product_status enum in Postgres.
-- The value will remain but is safe to ignore until the type is fully dropped.