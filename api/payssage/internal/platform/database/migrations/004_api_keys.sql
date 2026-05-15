-- +goose Up
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL,
    key_prefix TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_api_keys_workspace_id ON api_keys (workspace_id);
CREATE INDEX idx_api_keys_key_prefix ON api_keys (key_prefix);
-- +goose Down
DROP INDEX IF EXISTS idx_api_keys_key_prefix;
DROP INDEX IF EXISTS idx_api_keys_workspace_id;
DROP TABLE IF EXISTS api_keys;
