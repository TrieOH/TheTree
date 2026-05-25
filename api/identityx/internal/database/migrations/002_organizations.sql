-- +goose Up
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    owner_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE RESTRICT,

    name TEXT NOT NULL,
    slug TEXT NOT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uniq_organizations_slug
       UNIQUE (slug)
);

CREATE INDEX idx_organizations_created_at ON organizations (created_at);
CREATE INDEX idx_organizations_metadata_gin ON organizations USING GIN (metadata);

CREATE TABLE org_members (
    PRIMARY KEY (organization_id, actor_id),
    organization_id UUID NOT NULL REFERENCES organizations(id)
        ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,

    role TEXT NOT NULL,
    CONSTRAINT chk_org_members_role CHECK (role IN ('owner', 'admin', 'member')),

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_org_members_actor_id ON org_members (actor_id);
CREATE INDEX idx_org_members_role ON org_members (role);

-- +goose Down
DROP INDEX IF EXISTS idx_org_members_role;
DROP INDEX IF EXISTS idx_org_members_actor_id;
DROP TABLE IF EXISTS org_members;
DROP INDEX IF EXISTS idx_organizations_metadata_gin;
DROP INDEX IF EXISTS idx_organizations_created_at;
DROP TABLE IF EXISTS organizations;
