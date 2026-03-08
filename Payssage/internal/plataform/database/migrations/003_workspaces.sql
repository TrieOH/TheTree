-- +goose Up
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    sandbox BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_workspaces_user_name ON workspaces (user_id, name);
CREATE INDEX idx_workspaces_user_id ON workspaces (user_id);
-- +goose Down
DROP INDEX IF EXISTS idx_workspaces_user_id;
DROP INDEX IF EXISTS idx_workspaces_user_name;
DROP TABLE IF EXISTS workspaces;