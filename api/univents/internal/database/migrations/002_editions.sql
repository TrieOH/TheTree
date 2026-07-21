-- +goose Up
CREATE TABLE editions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    edition_name VARCHAR(256) NOT NULL,
    slug VARCHAR(32) NOT NULL,
    tagline VARCHAR(512),
    description TEXT,

    is_draft BOOLEAN NOT NULL DEFAULT TRUE,

    registration_opens_at TIMESTAMPTZ,

    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT chk_editions_dates_valid CHECK (ends_at > starts_at),
    CONSTRAINT chk_editions_registration_before_start CHECK (
        registration_opens_at IS NULL OR registration_opens_at <= starts_at
    ),

    location_name TEXT,
    location_address TEXT,

    logo_url TEXT,
    banner_url TEXT,
    contact_email VARCHAR(256),

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_edition_slug_unique
    ON editions(event_id, slug) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_edition_slug_unique;
DROP TABLE IF EXISTS editions;
