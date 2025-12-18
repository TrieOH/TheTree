-- +goose Up
-- Created at 2025-11-04T15:25:15-03:00

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE project_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_type TEXT NOT NULL DEFAULT 'project',
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT True,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP,
    UNIQUE (project_id, email)
);

ALTER TABLE user_sessions
    ADD COLUMN user_type TEXT NOT NULL;

ALTER TABLE users
    ADD COLUMN user_type TEXT NOT NULL DEFAULT 'client';

CREATE INDEX IF NOT EXISTS idx_project_users_project_id
    ON project_users(project_id);

-- +goose Down
DROP TABLE IF EXISTS project_users;
DROP INDEX IF EXISTS idx_project_users_project_id;

ALTER TABLE user_sessions
DROP COLUMN IF EXISTS user_type;
