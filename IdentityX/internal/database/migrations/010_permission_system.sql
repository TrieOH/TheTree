-- +goose Up

-- FIXME Add deprecated_at later for scope deprecation (soft_delete)
CREATE TABLE scopes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Hierarchy reference (null only for global scope)
    parent_id UUID REFERENCES scopes(id) ON DELETE RESTRICT,

    -- Internal hierarchy type (controlled by the IdP only)
    type TEXT NOT NULL CHECK (type IN ('global', 'project_root', 'project_scope')),

    -- Tenant boundary (NULL only for global)
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,

    -- Namespace inside a project (NULL for global + project_root)
    name TEXT,

    UNIQUE NULLS DISTINCT (type, name),

    -- Optional reference to a specific resource inside that namespace
    external_id TEXT,

    meta JSONB NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT scope_shape_check CHECK (
        -- 🌍 One true global scope (root of IdP scopes, isolated from projects)
        (type = 'global'
            AND project_id IS NULL
            AND parent_id IS NULL
            AND name IS NULL
            AND external_id IS NULL)

            OR

        -- 🏗 One root scope per project (root of its own hierarchy, isolated from IdP)
        (type = 'project_root'
            AND project_id IS NOT NULL
            AND parent_id IS NULL
            AND name IS NULL
            AND external_id IS NULL)

            OR

        -- 📦 Named scopes inside a project (hierarchical, arbitrary depth)
        (type = 'project_scope'
            AND project_id IS NOT NULL
            AND parent_id IS NOT NULL
            AND name IS NOT NULL)
        )
);

-- Only ONE global scope in the whole system
CREATE UNIQUE INDEX scopes_one_global
    ON scopes (type)
    WHERE type = 'global';

-- Only ONE project_root scope per project
CREATE UNIQUE INDEX scopes_one_project_root_per_project
    ON scopes (project_id)
    WHERE type = 'project_root';

-- Unique named scopes per parent (sibling uniqueness)
CREATE UNIQUE INDEX scopes_unique_siblings
    ON scopes (parent_id, name)
    WHERE external_id IS NULL;

-- Unique resource-bound scopes per parent (sibling uniqueness with external_id)
CREATE UNIQUE INDEX scopes_unique_resource_siblings
    ON scopes (parent_id, name, external_id)
    WHERE external_id IS NOT NULL;

-- Fast parent lookups for hierarchy traversal
CREATE INDEX idx_scopes_parent_id ON scopes(parent_id);

-- Index for project scope listings
CREATE INDEX idx_scopes_project_id ON scopes(project_id);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    object TEXT NOT NULL, -- plain object name, e.g. "document", "event"
    action TEXT NOT NULL, -- plain action name, e.g. "read", "write"
    meta JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (length(object) BETWEEN 1 AND 100),
    CHECK (length(action) BETWEEN 1 AND 100),
    CHECK (object ~ '^(\*|[a-zA-Z][a-zA-Z0-9_]*)$'),
    CHECK (action ~ '^(\*|[a-zA-Z][a-zA-Z0-9_]*)$'),

    UNIQUE NULLS NOT DISTINCT (project_id, object, action)
);

CREATE INDEX idx_permissions_object_action
    ON permissions(object, action);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- null = IdP role, not null = project role
    project_id UUID NULL REFERENCES projects(id) ON DELETE CASCADE,

    name VARCHAR(64) NOT NULL,
    description TEXT,

    meta JSONB NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- unique: (admin, null) once, (admin, proj-a) once, (admin, proj-b) once...
    UNIQUE NULLS NOT DISTINCT (name, project_id)
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,

    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE identity_roles (
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    scope_id UUID REFERENCES scopes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE NULLS NOT DISTINCT (identity_id, role_id, scope_id)
);

CREATE INDEX idx_identity_roles_identity_scope
    ON identity_roles(identity_id, scope_id);

CREATE INDEX idx_identity_roles_lookup
    ON identity_roles(identity_id, scope_id, role_id);

CREATE TABLE identity_permissions (
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    scope_id UUID REFERENCES scopes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE NULLS NOT DISTINCT (identity_id, permission_id, scope_id)
);

CREATE INDEX idx_identity_permissions_identity_scope
    ON identity_permissions(identity_id, scope_id);

CREATE INDEX idx_identity_permissions_lookup
    ON identity_permissions(identity_id, scope_id, permission_id);

CREATE TABLE permission_audit_log (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL REFERENCES identities(id),
    target_identity_id UUID REFERENCES identities(id),
    action TEXT NOT NULL, -- grant_role | revoke_role | grant_permission
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS permission_audit_log;
DROP INDEX IF EXISTS idx_identity_permissions_lookup;
DROP INDEX IF EXISTS idx_identity_permissions_identity_scope;
DROP TABLE IF EXISTS identity_permissions;
DROP INDEX IF EXISTS idx_identity_roles_lookup;
DROP INDEX IF EXISTS idx_identity_roles_identity_scope;
DROP TABLE IF EXISTS identity_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP INDEX IF EXISTS idx_permissions_object_action;
DROP TABLE IF EXISTS permissions;
DROP INDEX IF EXISTS scopes_unique_resource_siblings;
DROP INDEX IF EXISTS scopes_unique_siblings;
DROP INDEX IF EXISTS idx_scopes_project_id;
DROP INDEX IF EXISTS idx_scopes_parent_id;
DROP INDEX IF EXISTS scopes_one_project_root_per_project;
DROP INDEX IF EXISTS scopes_one_global;
DROP TABLE IF EXISTS scopes;
