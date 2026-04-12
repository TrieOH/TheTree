-- +goose Up
-- 010_create_api_keys.sql
CREATE TABLE api_keys (
    project_id UUID PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE api_keys;
