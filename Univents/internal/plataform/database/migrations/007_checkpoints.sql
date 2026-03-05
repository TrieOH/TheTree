-- +goose Up

-- Checkpoints (passive tracking - coffee breaks, zones, unticketed access)
CREATE TYPE checkpoint_type AS ENUM (
    'entry',        -- venue entry
    'zone',         -- specific area (vip, backstage)
    'amenity',      -- coffee, lunch, swag pickup
    'session',      -- unticketed session tracking
    'exit'          -- departure tracking
);

CREATE TYPE checkpoint_access AS ENUM (
    'open',         -- no ticket required, just count
    'ticket',       -- any valid ticket
    'staff_only'    -- manual override required
);

CREATE TABLE checkpoints (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,

    name VARCHAR(256) NOT NULL,
    type checkpoint_type NOT NULL,

    starts_at TIMESTAMPTZ NULL,
    ends_at TIMESTAMPTZ NULL,

    -- access rules
    access_mode checkpoint_access NOT NULL DEFAULT 'open',

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_checkpoints_edition ON checkpoints(edition_id, type, access_mode)
    WHERE deleted_at IS NULL;

-- Checkpoint scans (who entered where/when)
CREATE TABLE checkpoint_scans (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    checkpoint_id UUID NOT NULL REFERENCES checkpoints(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,

    -- access validation at time of scan
    access_granted BOOLEAN NOT NULL,
    denial_reason VARCHAR(64) NULL, -- 'no_ticket', 'not_open'

    -- what was used
    --user_ticket_id UUID NULL REFERENCES user_tickets(id),

    -- scan metadata
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    scanned_by UUID NULL, -- staff member, null for self-scan
    device_id VARCHAR(64) NULL, -- scanner device

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_checkpoint_scans_checkpoint ON checkpoint_scans(checkpoint_id, scanned_at);
CREATE INDEX idx_checkpoint_scans_user ON checkpoint_scans(user_id, scanned_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_checkpoint_scans_user;
DROP INDEX IF EXISTS idx_checkpoint_scans_checkpoint;
DROP TABLE IF EXISTS checkpoint_scans;
DROP INDEX IF EXISTS idx_checkpoints_edition;
DROP TABLE IF EXISTS checkpoints;
DROP TYPE IF EXISTS checkpoint_access;
DROP TYPE IF EXISTS checkpoint_type;
