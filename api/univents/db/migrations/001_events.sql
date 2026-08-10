-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,

    full_name VARCHAR(256) NOT NULL,
    acronym VARCHAR(32),
    slug VARCHAR(32) NOT NULL,
    description TEXT,

    style JSONB,

    status TEXT NOT NULL DEFAULT 'draft',
    CONSTRAINT chk_event_status_valid CHECK (
        status IN (
            'draft',
            'active',
            'discontinued'
        )
    ),

    payssage_seller_id UUID,
    payssage_public_key TEXT,
    CONSTRAINT chk_event_payments_public_key_requires_seller CHECK (
        payssage_public_key IS NULL OR payssage_seller_id IS NOT NULL
    ),

    logo_url TEXT,
    banner_url TEXT,
    contact_email VARCHAR(256),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_slug ON events(slug)
    WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_events_slug_unique ON events(slug)
    WHERE deleted_at IS NULL;

CREATE TABLE event_members (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    event_id   UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,
    role       TEXT NOT NULL DEFAULT 'staff',
    CONSTRAINT chk_event_members_role_valid CHECK (
        role IN ('owner', 'admin', 'staff')
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_one_event_member_per_event
    ON event_members(event_id, user_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_one_event_member_per_event;
DROP TABLE IF EXISTS event_members;
DROP INDEX IF EXISTS idx_events_slug_unique;
DROP INDEX IF EXISTS idx_events_slug;
DROP TABLE IF EXISTS events;
DROP EXTENSION IF EXISTS "pgcrypto";