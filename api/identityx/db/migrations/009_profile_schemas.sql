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

-- history of profile schema versions. Every upsert that changes a schema
-- appends a row here, so each version stays reproducible. project_id NULL
-- is the platform-wide singleton scope, mirroring project_profile_schemas.
CREATE TABLE profile_schema_versions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID REFERENCES projects(id)
        ON DELETE CASCADE, -- NULL = platform scope

    version INTEGER NOT NULL,
    schema  JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_profile_schema_version_per_scope
        UNIQUE (project_id, version)
);

-- the platform singleton scope (project_id NULL) holds a single version row
CREATE UNIQUE INDEX uniq_profile_schema_version_platform
    ON profile_schema_versions (project_id) WHERE project_id IS NULL;

CREATE INDEX idx_profile_schema_versions_project
    ON profile_schema_versions (project_id);

-- +goose Down
DROP INDEX IF EXISTS idx_profile_schema_versions_project;
DROP INDEX IF EXISTS uniq_profile_schema_version_platform;
DROP TABLE IF EXISTS profile_schema_versions;
DROP INDEX IF EXISTS uniq_project_profile_schema_project_id;
DROP TABLE IF EXISTS project_profile_schemas;
