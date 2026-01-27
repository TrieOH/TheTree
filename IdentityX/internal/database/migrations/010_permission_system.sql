-- +goose Up

-- FIXME Add deprecated_at later for scope deprecation (soft_delete)
CREATE TABLE scopes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- Internal hierarchy (controlled by the IdP only)
    type TEXT NOT NULL CHECK (type IN ('global', 'project_root', 'project_scope')),

    -- Tenant boundary (NULL only for global)
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,

    -- Namespace inside a project (NULL for global + project_root)
    name TEXT,

    -- Optional reference to a specific resource inside that namespace
    external_id TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT scope_shape_check CHECK (

        -- 🌍 One true global scope
        (type = 'global'
            AND project_id IS NULL
            AND name IS NULL
            AND external_id IS NULL)

            OR

        -- 🏗 One root scope per project
        (type = 'project_root'
            AND project_id IS NOT NULL
            AND name IS NULL
            AND external_id IS NULL)

            OR

        -- 📦 Named scopes inside a project
        (type = 'project_scope'
            AND project_id IS NOT NULL
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

-- Unique named scopes per project
CREATE UNIQUE INDEX scopes_unique_project_named_scopes
    ON scopes (project_id, name)
    WHERE type = 'project_scope' AND external_id IS NULL;

-- Unique resource-bound scopes per project
CREATE UNIQUE INDEX scopes_unique_project_resource_scopes
    ON scopes (project_id, name, external_id)
    WHERE type = 'project_scope' AND external_id IS NOT NULL;


--- object: object:specifier (/ path)
--- object: object:specifier/object:specifier
--- having no specifier is the same as specifying *

--- action: domain:verb
--- or action: verb
--- action: attendance:mark
--- action: edit
--- actions should not have specifiers

--- Examples
--- object: event:123/activity:456
--- action: attendance:mark
--- User X Can mark attendance on activity 456 on event 123

--- object: event:123/activity:*
--- action: attendance:mark
--- User X Can mark attendance on all activities on event 123

--- object: event:123/activity
--- action: create
--- User X Can create activities on event 123

--- forbidden edit requires either all or specific
--- object: event:123/activity
--- action: edit
--- User X Can edit activities on event 123

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    object TEXT NOT NULL, -- e.g. "event:*", "event:123", "event:123/activity"
    action TEXT NOT NULL, -- e.g. "create", "edit", "delete", "attendance:mark"
    conditions JSONB,     -- optional ABAC rules
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (length(object) BETWEEN 1 AND 255),
    CHECK (length(action) BETWEEN 1 AND 100),
    CHECK (object ~ '^[a-zA-Z0-9:_/*-]+$'),
    CHECK (action ~ '^[a-zA-Z0-9:_*-]+$'),

    UNIQUE NULLS NOT DISTINCT (project_id, object, action)
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    -- null = IdP role, not null = project role
    project_id UUID NULL REFERENCES projects(id) ON DELETE CASCADE,

    name VARCHAR(64) NOT NULL,
    description TEXT,

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
    scope_id UUID NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (identity_id, role_id, scope_id)
);

CREATE TABLE identity_permissions (
    identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    scope_id UUID NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (identity_id, permission_id, scope_id)
);

CREATE TABLE permission_audit_log (
    id UUID PRIMARY KEY,
    identity_id UUID NOT NULL REFERENCES identities(id),
    target_identity_id UUID REFERENCES identities(id),
    action TEXT NOT NULL, -- grant_role | revoke_role | grant_permission
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_identity_roles_identity_scope
    ON identity_roles (identity_id, scope_id);

CREATE INDEX idx_identity_permissions_identity_scope
    ON identity_permissions (identity_id, scope_id);

CREATE INDEX idx_permissions_object_action
    ON permissions (object, action);

-- +goose Down
DROP INDEX IF EXISTS idx_permissions_object_action;
DROP INDEX IF EXISTS idx_identity_permissions_identity_scope;
DROP INDEX IF EXISTS idx_identity_roles_identity_scope;
DROP TABLE IF EXISTS permission_audit_log;
DROP TABLE IF EXISTS identity_permissions;
DROP TABLE IF EXISTS identity_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP INDEX IF EXISTS permissions_project_unique;
DROP INDEX IF EXISTS permissions_idp_unique;
DROP TABLE IF EXISTS permissions;
DROP INDEX IF EXISTS scopes_unique_project_resource_scopes;
DROP INDEX IF EXISTS scopes_unique_project_named_scopes;
DROP INDEX IF EXISTS scopes_one_project_root_per_project;
DROP INDEX IF EXISTS scopes_one_global;
DROP TABLE IF EXISTS scopes;
