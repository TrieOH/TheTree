-- +goose Up
-- Created at 2025-11-04T15:25:15-03:00

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_name VARCHAR(255) NOT NULL DEFAULT 'Unnamed Project',
    owner_id UUID NOT NULL REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT True,
    pub_key TEXT NOT NULL,
    priv_key BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE user_sessions
    ADD COLUMN project_id UUID REFERENCES projects(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS projects;
DROP EXTENSION IF EXISTS "pgcrypto";

ALTER TABLE user_sessions
DROP COLUMN IF EXISTS project_id;
