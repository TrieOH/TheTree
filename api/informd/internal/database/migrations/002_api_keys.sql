-- +goose Up
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    key_prefix TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ NULL,

    CONSTRAINT uniq_name_per_user UNIQUE (owner_id, name)
);

CREATE INDEX idx_api_keys_owner_id ON api_keys (owner_id);
CREATE INDEX idx_api_keys_key_prefix ON api_keys (key_prefix);

-- +goose Down
DROP INDEX IF EXISTS idx_api_keys_key_prefix;
DROP INDEX IF EXISTS idx_api_keys_owner_id;
DROP INDEX IF EXISTS uniq_idx_api_keys_name_project;
DROP TABLE IF EXISTS api_keys;
