-- +goose Up
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    organization_id UUID REFERENCES organizations(id)
        ON DELETE CASCADE,

    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    CONSTRAINT uniq_projects_slug UNIQUE (slug),

    domain TEXT,
    CONSTRAINT uniq_verified_domain UNIQUE(domain),
    domain_verified_at TIMESTAMPTZ,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_projects_organization_id ON projects(organization_id);
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_metadata_gin ON projects USING GIN (metadata);

CREATE TABLE project_domain_challenges(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,

    domain TEXT NOT NULL,

    -- HIGH_ENTROPY_HASH
    token TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,

    -- NULL = pending
    verified_at TIMESTAMPTZ
);

CREATE INDEX idx_project_domain_challenges ON project_domain_challenges(project_id);
CREATE INDEX idx_domain_project_domain_challenges ON project_domain_challenges(domain);

CREATE TABLE project_members (
    PRIMARY KEY (project_id, actor_id),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,

    role TEXT NOT NULL,
    CONSTRAINT chk_project_members_role CHECK (
        role IN (
             'owner',
             'admin',
             'developer',
             'analyst',
             'support'
        )
    ),

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_project_members_actor_id ON project_members (actor_id);
CREATE INDEX idx_project_members_role ON project_members (role);

CREATE TABLE project_oauth_providers (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id)
        ON DELETE CASCADE,

    provider TEXT NOT NULL,
    CONSTRAINT chk_project_oauth_providers_provider CHECK (
        provider IN ('google', 'github')
    ),
    CONSTRAINT uniq_project_oauth_provider UNIQUE (project_id, provider),

    client_id TEXT NOT NULL,
    encrypted_client_secret TEXT NOT NULL,

    scopes TEXT[] NOT NULL DEFAULT '{}',

    enabled BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_project_oauth_providers_project_id
    ON project_oauth_providers (project_id);

-- +goose Down
DROP INDEX IF EXISTS idx_project_oauth_providers_project_id;
DROP TABLE IF EXISTS project_oauth_providers;
DROP INDEX IF EXISTS idx_project_members_role;
DROP INDEX IF EXISTS idx_project_members_actor_id;
DROP TABLE IF EXISTS project_members;
DROP INDEX IF EXISTS idx_domain_project_domain_challenges;
DROP INDEX IF EXISTS idx_project_domain_challenges;
DROP TABLE IF EXISTS project_domain_challenges;
DROP INDEX IF EXISTS idx_projects_metadata_gin;
DROP INDEX IF EXISTS idx_projects_created_at;
DROP INDEX IF EXISTS idx_projects_organization_id;
DROP TABLE IF EXISTS projects;
