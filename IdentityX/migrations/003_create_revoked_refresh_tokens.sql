-- +goose Up
-- Created at 2025-10-12T22:54:36-03:00

CREATE TABLE revoked_refresh_tokens (
    token_id UUID PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_revoked_refresh_tokens_token_id
    ON revoked_refresh_tokens(token_id);

-- +goose Down
DROP INDEX IF EXISTS idx_revoked_refresh_tokens_token_id;
DROP TABLE IF EXISTS revoked_refresh_tokens;
