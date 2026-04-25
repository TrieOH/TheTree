-- +goose Up
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uniq_idx_projects_owner_id_name ON projects (owner_id, name);
CREATE INDEX idx_projects_owner_id ON projects (owner_id);
-- +goose Down
DROP INDEX IF EXISTS idx_projects_owner_id;
DROP INDEX IF EXISTS uniq_idx_projects_owner_id_name;
DROP TABLE IF EXISTS projects;