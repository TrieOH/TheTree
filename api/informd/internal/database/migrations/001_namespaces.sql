-- +goose Up
CREATE TABLE namespaces (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uniq_namespace_name_per_user UNIQUE (owner_id, name)
);

CREATE INDEX idx_namespaces_owner_id ON namespaces (owner_id);

CREATE TABLE namespace_members (
    user_id UUID NOT NULL,
    namespace_id UUID NOT NULL REFERENCES namespaces(id)
        ON DELETE CASCADE,

    role VARCHAR(32) NOT NULL,
    CONSTRAINT chk_valid_namespace_member_role CHECK (
        role in ('viewer', 'editor', 'admin', 'owner')
    ),

    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_by UUID NOT NULL,

    PRIMARY KEY (user_id, namespace_id)
);

CREATE INDEX idx_namespace_members_role ON namespace_members (role);

-- +goose Down
DROP INDEX IF EXISTS idx_namespace_members_role;
DROP TABLE IF EXISTS namespace_members;
DROP INDEX IF EXISTS idx_namespaces_owner_id;
DROP TABLE IF EXISTS namespaces;