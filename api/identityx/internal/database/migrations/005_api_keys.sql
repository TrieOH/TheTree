-- +goose Up
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id)
        ON DELETE SET NULL,

    name TEXT NOT NULL,

    key_prefix TEXT NOT NULL,
    key_hash TEXT NOT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uniq_idx_api_keys_key_prefix ON api_keys (key_prefix);
CREATE INDEX idx_api_keys_actor_id ON api_keys (actor_id);
CREATE INDEX idx_api_keys_project_id ON api_keys (project_id);
CREATE INDEX idx_api_keys_revoked_at ON api_keys (revoked_at);
CREATE INDEX idx_api_keys_expires_at ON api_keys (expires_at);
-- +goose Down
DROP INDEX IF EXISTS idx_api_keys_expires_at;
DROP INDEX IF EXISTS idx_api_keys_revoked_at;
DROP INDEX IF EXISTS idx_api_keys_project_id;
DROP INDEX IF EXISTS idx_api_keys_actor_id;
DROP INDEX IF EXISTS uniq_idx_api_keys_key_prefix;
DROP TABLE IF EXISTS api_keys;
