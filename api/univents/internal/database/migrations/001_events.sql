-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,

    full_name VARCHAR(256) NOT NULL,
    acronym VARCHAR(32),
    slug VARCHAR(32) NOT NULL,
    description TEXT NULL,

    status TEXT NOT NULL DEFAULT 'draft',
    CONSTRAINT chk_event_status_valid CHECK (
        status IN (
            'draft',
            'active',
            'discontinued'
        )
    ),

    payssage_seller_id UUID,
    payssage_wallet_id UUID,
    CONSTRAINT chk_event_payments_config_complete CHECK (
        (payssage_seller_id IS NULL) = (payssage_wallet_id IS NULL)
    ),

    logo_url TEXT NULL,
    banner_url TEXT NULL,
    contact_email VARCHAR(256) NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_slug ON events(slug)
    WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_events_slug_unique ON events(slug)
    WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_events_slug_unique;
DROP INDEX IF EXISTS idx_events_slug;
DROP TABLE IF EXISTS events;
DROP EXTENSION IF EXISTS "pgcrypto";