-- +goose Up
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    subject_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,

    name TEXT NOT NULL,

    display_prefix TEXT NOT NULL,
    key_hash BYTEA NOT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,

    created_by UUID NOT NULL REFERENCES actors(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uniq_idx_api_keys_display_prefix ON api_keys (display_prefix);
CREATE INDEX idx_api_keys_subject_id ON api_keys (subject_id);
CREATE INDEX idx_api_keys_revoked_at ON api_keys (revoked_at);
CREATE INDEX idx_api_keys_expires_at ON api_keys (expires_at);
-- +goose Down
DROP INDEX IF EXISTS idx_api_keys_expires_at;
DROP INDEX IF EXISTS idx_api_keys_revoked_at;
DROP INDEX IF EXISTS idx_api_keys_subject_id;
DROP INDEX IF EXISTS uniq_idx_api_keys_display_prefix;
DROP TABLE IF EXISTS api_keys;
