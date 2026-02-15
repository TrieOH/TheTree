-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- we are using pg18 so UUIDv7() is fine

CREATE TYPE event_status AS ENUM (
    'draft',        -- not visible
    'active',       -- visible, accepting editions
    'archived',     -- no new editions
    'discontinued'  -- permanently ended
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    organization_id UUID NULL,
    goauth_scope_id UUID NOT NULL,

    -- identity
    name VARCHAR(256) NOT NULL,
    acronym VARCHAR(32) NULL,
    slug VARCHAR(32) NOT NULL,
    tagline VARCHAR(512) NULL,
    description TEXT NULL,

    -- classification
    is_series BOOLEAN NOT NULL DEFAULT FALSE, -- if false limits to one edition
    editions_count INT NOT NULL DEFAULT 0,

    CONSTRAINT chk_series_requires_single_edition
        CHECK (is_series = TRUE OR editions_count <= 1),

    -- branding
    logo_url TEXT NULL,
    banner_url TEXT NULL,

    -- images
    has_gallery BOOLEAN NOT NULL DEFAULT FALSE,
    gallery_urls TEXT[] NULL,

    -- contact
    contact_email VARCHAR(256) NULL,
    social_links JSONB NULL,  -- {twitter: "...", linkedin: "..."}

    -- state
    status event_status NOT NULL DEFAULT 'draft',

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_events_org_status ON events(organization_id, status)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_events_slug ON events(slug)
    WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_events_slug_unique ON events(slug)
    WHERE deleted_at IS NULL;

CREATE TYPE event_audit_action AS ENUM (
    -- lifecycle
    'created',
    'edited',
    'activated',
    'archived',
    'discontinued',
    'deleted',
    'restored',

    -- branding
    'logo_updated',
    'banner_updated',
    'gallery_updated',

    -- metadata
    'slug_changed',
    'contact_updated',
    'social_links_updated',

    -- editions
    'edition_added',
    'edition_removed',

    -- access
    'scope_changed',
    'ownership_transferred'
);

CREATE TYPE event_actor_type AS ENUM (
    'owner',
    'admin',
    'staff',
    'system'
);

CREATE TABLE event_audit (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    actor_type event_actor_type NOT NULL,
    actor_id UUID NULL,

    action event_audit_action NOT NULL,

    -- state tracking
    from_status event_status NULL,
    to_status event_status NULL,

    -- change context
    metadata JSONB NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_event_audit_event_lookup ON event_audit(event_id, created_at DESC);
CREATE INDEX idx_event_audit_action ON event_audit(action, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_event_audit_action;
DROP INDEX IF EXISTS idx_event_audit_event_lookup;
DROP TABLE IF EXISTS event_audit;
DROP TYPE IF EXISTS event_actor_type;
DROP TYPE IF EXISTS event_audit_action;
DROP INDEX IF EXISTS idx_events_slug_unique;
DROP INDEX IF EXISTS idx_events_slug;
DROP INDEX IF EXISTS idx_events_org_status;
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;
DROP EXTENSION IF EXISTS "pgcrypto";