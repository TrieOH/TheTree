-- +goose Up
CREATE TABLE project_profile_schemas (
    project_id UUID REFERENCES projects(id)
        ON DELETE CASCADE,

    schema JSONB NOT NULL,

    version  INTEGER NOT NULL DEFAULT 1,
    active   BOOLEAN NOT NULL DEFAULT true,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- only one schema per project, including the NULL→platform singleton
CREATE UNIQUE INDEX uniq_project_profile_schema_project_id
    ON project_profile_schemas (project_id) NULLS NOT DISTINCT;

-- +goose Down
DROP INDEX IF EXISTS uniq_project_profile_schema_project_id;
DROP TABLE IF EXISTS project_profile_schemas;
