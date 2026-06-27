-- +goose Up
CREATE TABLE organizations (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id   UUID NOT NULL,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uniq_organizations_slug UNIQUE (slug)
);

CREATE TABLE org_members (
    PRIMARY KEY (organization_id, member_id),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    member_id       UUID NOT NULL,
    role            TEXT NOT NULL,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_org_members_role CHECK (role IN ('owner', 'admin', 'member'))
);

CREATE INDEX idx_org_members_member_id ON org_members (member_id);

-- +goose Down
DROP INDEX IF EXISTS idx_org_members_member_id;
DROP TABLE IF EXISTS org_members;
DROP TABLE IF EXISTS organizations;