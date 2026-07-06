-- +goose Up
CREATE TABLE capabilities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID REFERENCES projects(id)
        ON DELETE CASCADE,

    resource TEXT NOT NULL,
    action TEXT NOT NULL,

    created_by UUID NOT NULL REFERENCES actors(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_capability_per_scope UNIQUE NULLS NOT DISTINCT (project_id, resource, action)
);

CREATE INDEX idx_capabilities_project ON capabilities(project_id);

CREATE TABLE actor_capabilities (
    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES capabilities(id)
        ON DELETE CASCADE,

    assigned_by UUID REFERENCES actors(id),
    assigned_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(actor_id, capability_id)
);

CREATE TABLE api_key_capabilities (
    api_key_id UUID NOT NULL REFERENCES api_keys(id)
        ON DELETE CASCADE,
    capability_id UUID NOT NULL REFERENCES capabilities(id)
        ON DELETE CASCADE,

    assigned_by UUID REFERENCES actors(id),
    assigned_at TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(api_key_id, capability_id)
);

-- +goose Down
DROP TABLE IF EXISTS api_key_capabilities;
DROP TABLE IF EXISTS actor_capabilities;
DROP INDEX IF EXISTS idx_capabilities_project;
DROP TABLE IF EXISTS capabilities;
