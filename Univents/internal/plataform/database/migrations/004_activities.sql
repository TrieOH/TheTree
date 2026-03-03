-- +goose Up

CREATE TYPE activity_status AS ENUM (
    'draft',        -- hidden, editing
    'published',    -- visible, open for registration
    'ongoing',      -- currently happening, check-in open
    'completed',    -- ended, read-only
    'canceled'     -- activity canceled, not just attendance
);

CREATE TYPE difficulty_level AS ENUM (
    'no_prerequisites', -- No prior knowledge or skills required
    'beginner',         -- Basic familiarity helpful but not required
    'intermediate',     -- Some prior experience recommended
    'advanced',         -- Solid foundation in subject required
    'expert'            -- Deep expertise assumed, fast-paced
);

CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    title VARCHAR(256) NOT NULL,
    description TEXT NULL,
    status activity_status NOT NULL DEFAULT 'draft',
    location VARCHAR(256) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,

    CHECK (ends_at > starts_at),

    -- images
    -- has_activity_banner BOOLEAN NOT NULL DEFAULT FALSE,
    -- activity_banner_url VARCHAR(1024) NULL,
    -- has_activity_photo BOOLEAN NOT NULL DEFAULT FALSE,
    -- activity_photo_url VARCHAR(1024) NULL,
    -- has_gallery BOOLEAN NOT NULL DEFAULT FALSE,
    -- gallery_image_urls VARCHAR(1024)[] NULL,

    -- presenter
    -- has_presenter BOOLEAN NOT NULL DEFAULT FALSE,
    -- presenter_id UUID NULL,
    presenter_name VARCHAR(256) NULL,
    -- has_presenter_photo BOOLEAN NOT NULL DEFAULT FALSE,
    -- presenter_photo_url VARCHAR(1024) NULL,

    -- tokens
    --is_blocked_by_tokens BOOLEAN NOT NULL DEFAULT FALSE,
    token_cost INT NOT NULL DEFAULT 0,

    -- paid entry
    -- has_paid_entry BOOLEAN NOT NULL DEFAULT FALSE, -- if an event wants to sell access to an activity without the participant having to buy access to the event, event participants get it for free no matter the price unless its blocked by tokens
    -- price_cents INT NOT NULL DEFAULT 0,

    -- capacity
    has_capacity BOOLEAN NOT NULL DEFAULT FALSE,
    capacity INT NOT NULL DEFAULT 0,
    remaining_capacity INT NOT NULL DEFAULT 0,

    CHECK (has_capacity = FALSE OR capacity > 0),
    CHECK (capacity >= remaining_capacity),

    -- waitlist (depends on capacity)
    -- has_waitlist BOOLEAN NOT NULL DEFAULT FALSE,
    -- has_waitlist_capacity BOOLEAN NOT NULL DEFAULT FALSE,
    -- waitlist_capacity INT NOT NULL DEFAULT 0,
    -- waitlist_remaining_capacity INT NOT NULL DEFAULT 0,

    -- CHECK (has_waitlist = FALSE OR waitlist_capacity > 0),
    -- CHECK (waitlist_capacity >= waitlist_remaining_capacity),

    -- attendance
    --has_attendance BOOLEAN NOT NULL DEFAULT FALSE,
    --minimum_duration_for_attendance_minutes INT NOT NULL DEFAULT 0,

    -- check-in
    --check_in_closes_at TIMESTAMPTZ NULL,                        -- explicit time override
    --check_in_grace_before_start_minutes INT NOT NULL DEFAULT 0, -- time before activity start that check in is allowed

    -- checkout & duration
    --check_out_grace_after_end_minutes INT NOT NULL DEFAULT 0, -- time after activity ends that check out is allowed, in case of delays

    -- ratings
    -- has_ratings BOOLEAN NOT NULL DEFAULT FALSE,
    -- prompt_rating_at_exit BOOLEAN NOT NULL DEFAULT FALSE,

    -- materials
    -- has_materials BOOLEAN NOT NULL DEFAULT FALSE,
    -- materials JSONB NULL,

    -- difficulty
    difficulty difficulty_level NULL,

    -- requirements
    -- has_requirements BOOLEAN NOT NULL DEFAULT FALSE,
    -- requirements TEXT[] NULL,  -- ['React', 'Docker', 'Laptop']

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_activities_edition_lookup ON activities(edition_id, status, starts_at);

CREATE TABLE activity_interest_list (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,

    notified_when_open BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    UNIQUE(activity_id, user_id)
);

CREATE INDEX idx_activity_interest_activity ON activity_interest_list(activity_id, created_at);
CREATE INDEX idx_activity_interest_user ON activity_interest_list(user_id);

CREATE TYPE attendance_status AS ENUM (
    'registered',      -- signed up, not yet checked in
    'waitlisted',      -- on waitlist, not confirmed
    'promoted',        -- moved from waitlist to registered
    'checked_in',      -- arrived, currently attending
    'checked_out',     -- left, pending validation
    'completed',       -- met all requirements (duration, checkout)
    'partial',         -- attended but didn't meet min duration
    'no_show',         -- never checked in
    'cancelled'        -- user cancelled registration
);

CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    registration_times INT NOT NULL DEFAULT 1,

    status attendance_status NOT NULL,
    checked_in_at TIMESTAMPTZ NULL,
    checked_out_at TIMESTAMPTZ NULL,
    cancelled_at TIMESTAMPTZ NULL,

    rating INT NULL CHECK (rating >= 1 AND rating <= 5),
    rating_comment TEXT NULL,
    rated_at TIMESTAMPTZ NULL,

    is_paid_entry BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,

    UNIQUE(activity_id, user_id, registration_times)
);

CREATE TABLE waitlist_entries (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    position DECIMAL(10,4) NOT NULL DEFAULT 0, -- order in line
    promoted_at TIMESTAMPTZ NULL,              -- when moved to confirmed
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(activity_id, user_id)
);

CREATE INDEX idx_waitlist_position ON waitlist_entries(activity_id, position);

CREATE TYPE audit_action AS ENUM (
    -- attendance lifecycle
    'registered',
    'waitlisted',
    'promoted',
    'checked_in',
    'checked_out',
    'completed',
    'marked_partial',
    'marked_no_show',
    'registration_cancelled',
    'rating_submitted',

    -- activity management
    'created',
    'edited',
    'published',
    'started',
    'ended',
    'activity_cancelled',
    'activity_reopened',
    'deleted',

    -- capacity/waitlist
    'capacity_increased',
    'capacity_decreased',
    'waitlist_promoted'
);

CREATE TYPE actor_type AS ENUM (
    'participant',
    'staff',
    'admin',
    'owner',
    'system'
);

CREATE TABLE activity_audit (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- polymorphic target (activity, attendance_record, etc.)
    table_name VARCHAR(64) NOT NULL,
    record_id UUID NOT NULL,

    -- who did it
    actor_type actor_type NOT NULL DEFAULT 'participant',
    actor_id UUID NULL,  -- null for system

    action audit_action NOT NULL,

    -- state change (optional)
    state_type VARCHAR(64) NULL,
    from_status VARCHAR(32) NULL,
    to_status VARCHAR(32) NULL,

    -- change details
    metadata JSONB NULL,        -- context

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_lookup ON activity_audit(table_name, record_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_audit_lookup;
DROP TABLE IF EXISTS activity_audit;
DROP TYPE IF EXISTS actor_type;
DROP TYPE IF EXISTS audit_action;
DROP INDEX IF EXISTS idx_waitlist_position;
DROP TABLE IF EXISTS waitlist_entries;
DROP TABLE IF EXISTS attendance_records;
DROP TYPE IF EXISTS attendance_status;
DROP INDEX IF EXISTS idx_activity_interest_user;
DROP INDEX IF EXISTS idx_activity_interest_activity;
DROP INDEX IF EXISTS idx_activities_edition_lookup;
DROP TABLE IF EXISTS activity_interest_list;
DROP TABLE IF EXISTS activities;
DROP TYPE IF EXISTS difficulty_level;
DROP TYPE IF EXISTS activity_status;
