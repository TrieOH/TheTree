-- +goose Up
-- Created at 2025-10-12T22:54:36-03:00

CREATE TABLE refresh_blacklist (
    token_id UUID PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS refresh_blacklists;
