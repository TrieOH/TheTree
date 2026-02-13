-- +goose Up

CREATE TYPE edition_type AS ENUM (
    'year',         -- 2024, 2025
    'season',       -- Spring 2024, Fall 2025
    'number',       -- Edition 1, Edition 42
    'ordinal',      -- 1st Edition, 42nd Edition
    'custom'        -- "Anniversary Special", "COVID Edition"
);

CREATE TYPE edition_status AS ENUM (
    'draft',        -- planning, not visible
    'announced',    -- visible, registration not open
    'open',         -- registration open
    'ongoing',      -- currently happening
    'completed',    -- ended
    'cancelled',    -- edition cancelled
    'postponed'     -- rescheduled to new dates
);

CREATE TYPE edition_monetary_type AS ENUM (
    'free', -- no paid activities and no paid tickets
    'paid', -- all paid tickets
    'mixed' -- at least one free ticket and activity
);

CREATE TABLE editions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    event_id UUID NOT NULL REFERENCES events(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    -- naming
    type edition_type NOT NULL,
    edition_name VARCHAR(256) NOT NULL,
    name_template VARCHAR(512) NULL,
    tagline VARCHAR(512) NULL,
    description TEXT NULL,

    -- status & visibility
    status edition_status NOT NULL DEFAULT 'draft',
    visible_from TIMESTAMPTZ NULL,
    registration_opens_at TIMESTAMPTZ NULL,
    registration_closes_at TIMESTAMPTZ NULL,

    CHECK (registration_closes_at > registration_opens_at),

    -- tickets
    monetary_type edition_monetary_type NOT NULL,
    has_tickets BOOLEAN NOT NULL DEFAULT FALSE, -- if false there is just a register button and makes the event free, if true capacity is defined by ticket amount
    has_ticket_tiers BOOLEAN NOT NULL DEFAULT FALSE, -- allows tickets to give access to different things, normally all tickets give access to the same things
    CONSTRAINT chk_tiers_require_tickets
        CHECK (has_ticket_tiers = FALSE OR has_tickets = TRUE),
    has_capacity BOOLEAN NOT NULL DEFAULT FALSE, -- only if it has_tickets is false, then you can define max participants
    capacity INT NOT NULL DEFAULT 0,
    remaining_capacity INT NOT NULL DEFAULT 0,

    CHECK (has_capacity = FALSE OR capacity > 0),
    CHECK (capacity >= remaining_capacity),

    -- products
    has_products BOOLEAN NOT NULL DEFAULT FALSE,
    has_bundles BOOLEAN NOT NULL DEFAULT FALSE,
    has_tokens BOOLEAN NOT NULL DEFAULT FALSE, -- tokens are an internal currency that events can use to block stuff, like block an activity behind 2 tokens, forces a token product on the store
    max_tokens_per_user INT NOT NULL DEFAULT 1,
    CONSTRAINT chk_tokens_config_requires_enabled
        CHECK (has_tokens = TRUE OR max_tokens_per_user = 1),

    -- dates
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',

    CHECK (ends_at > starts_at),

    -- location
    location_name VARCHAR(256) NOT NULL,
    location_address TEXT NOT NULL,

    -- branding
    logo_url TEXT NULL,
    banner_url TEXT NULL,

    -- media
    has_gallery BOOLEAN NOT NULL DEFAULT FALSE,
    gallery_urls TEXT[] NULL,

    -- schedule
    has_schedule BOOLEAN NOT NULL DEFAULT FALSE,
    schedule_id UUID NULL,

    -- activities
    has_activities BOOLEAN NOT NULL DEFAULT FALSE,
    has_activity_interest_marking BOOLEAN NOT NULL DEFAULT FALSE,
    has_paid_activities BOOLEAN NOT NULL DEFAULT FALSE,

    -- interest options
    has_interest_list BOOLEAN NOT NULL DEFAULT FALSE,
    interest_list_opens_at TIMESTAMPTZ NULL,

    -- checkout, if the event requires a checkout to mark the user as attended
    -- if single day, checkout at the exit at days end, if multiday checkout at exit at last days end
    has_checkout BOOLEAN NOT NULL DEFAULT FALSE,

    -- contact
    contact_email VARCHAR(256) NULL,
    contact_phone VARCHAR(32) NULL,
    organizer_name VARCHAR(256) NULL,

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_editions_event_status ON editions(event_id, status, starts_at DESC);
CREATE INDEX idx_editions_visible ON editions(visible_from, status)
    WHERE status IN ('announced', 'open', 'ongoing');

CREATE TABLE edition_interest_list (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,

    notified_when_open BOOLEAN NOT NULL DEFAULT FALSE,  -- auto-notified when registration opens
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    UNIQUE(edition_id, user_id)
);

CREATE INDEX idx_interest_edition ON edition_interest_list(edition_id, created_at);
CREATE INDEX idx_interest_user ON edition_interest_list(user_id);

CREATE TYPE edition_registration_status AS ENUM (
    'registered',      -- signed up, not yet checked in
    'checked_in',      -- arrived, currently attending
    'attended',        -- if the user checked out the event in the end or attended once
    'partial',         -- attended but didn't get to the end
    'no_show',         -- never checked in
    'cancelled'        -- user cancelled registration
);

CREATE TABLE edition_registrations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,

    status edition_registration_status NOT NULL DEFAULT 'registered',

    registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at TIMESTAMPTZ NULL,
    cancellation_reason TEXT NULL,

    checked_in_at TIMESTAMPTZ NULL,
    checked_in_by UUID NULL,
    checked_out_at TIMESTAMPTZ NULL,
    checked_out_by UUID NULL,

    -- metadata
    notes TEXT NULL, -- admin notes

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    UNIQUE(edition_id, user_id)
);

CREATE INDEX idx_registrations_edition_status ON edition_registrations(edition_id, status, registered_at);
CREATE INDEX idx_registrations_user ON edition_registrations(user_id, status);

-- +goose Down
DROP INDEX IF EXISTS idx_registrations_user;
DROP INDEX IF EXISTS idx_registrations_edition_status;
DROP TABLE IF EXISTS edition_registrations;
DROP TYPE IF EXISTS edition_registration_status;
DROP INDEX IF EXISTS idx_interest_user;
DROP INDEX IF EXISTS idx_interest_edition;
DROP TABLE IF EXISTS edition_interest_list;
DROP INDEX IF EXISTS idx_editions_visible;
DROP INDEX IF EXISTS idx_editions_event_status;
DROP TABLE IF EXISTS editions;
DROP TYPE IF EXISTS edition_monetary_type;
DROP TYPE IF EXISTS edition_status;
DROP TYPE IF EXISTS edition_type;