-- +goose Up
-- Created at 2025-11-04T15:25:15-03:00

CREATE TABLE project_users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_type TEXT NOT NULL DEFAULT 'project',
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT True,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    UNIQUE (project_id, email)
);

CREATE INDEX IF NOT EXISTS idx_project_users_project_id
    ON project_users(project_id);

CREATE INDEX IF NOT EXISTS idx_project_users_user_type
    ON project_users(user_type);

-- +goose Down
DROP INDEX IF EXISTS idx_project_users_project_id;
DROP INDEX IF EXISTS idx_project_users_user_type;
DROP TABLE IF EXISTS project_users;
