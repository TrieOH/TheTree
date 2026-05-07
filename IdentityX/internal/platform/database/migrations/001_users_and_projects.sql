-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_name VARCHAR(255) NOT NULL DEFAULT 'Unnamed Project',
    domain VARCHAR(1024) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT True,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_id
    ON projects(id);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_type TEXT NOT NULL DEFAULT 'client',
    project_id UUID NULL REFERENCES projects(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_verified BOOL NOT NULL DEFAULT False,
    verified_at TIMESTAMPTZ NULL,

    CONSTRAINT chk_valid_user_type CHECK (user_type in ('client', 'project'))
);

CREATE INDEX IF NOT EXISTS idx_users_id
    ON users(id);

CREATE INDEX IF NOT EXISTS idx_users_project_id
    ON users(project_id);

CREATE INDEX IF NOT EXISTS idx_users_user_type
    ON users(user_type);

ALTER TABLE projects ADD COLUMN owner_id UUID NOT NULL REFERENCES users(id)
    ON DELETE CASCADE;

CREATE UNIQUE INDEX IF NOT EXISTS one_email_for_client
    ON users(email)
    WHERE project_id IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS one_email_per_project
    ON users(project_id, email)
    WHERE project_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS one_email_per_project;
DROP INDEX IF EXISTS one_email_for_client;
DROP INDEX IF EXISTS idx_users_user_type;
DROP INDEX IF EXISTS idx_users_project_id;
DROP INDEX IF EXISTS idx_users_id;
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_projects_id;
DROP TABLE IF EXISTS projects;
DROP EXTENSION IF EXISTS "pgcrypto";