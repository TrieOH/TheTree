-- +goose Up
-- Created at 2025-10-17T10:24:48-03:00

CREATE TABLE sessions (
    session_id UUID PRIMARY KEY DEFAULT uuidv7(),
    family_id UUID NOT NULL DEFAULT uuidv7(),
    token_id UUID UNIQUE NOT NULL DEFAULT uuidv7(),
    issued_at TIMESTAMPTZ NOT NULL,
    user_agent TEXT NOT NULL,
    user_ip VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (revoked_at IS NULL OR revoked_at >= issued_at)
);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at
    ON sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_sessions_revoked_at
    ON sessions(revoked_at)
    WHERE revoked_at IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_sessions_revoked_at;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP TABLE IF EXISTS sessions;
