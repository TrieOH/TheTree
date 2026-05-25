-- +goose Up
CREATE TABLE platform_roles (
    actor_id UUID PRIMARY KEY REFERENCES actors(id)
        ON DELETE CASCADE,

    role TEXT NOT NULL,
    CONSTRAINT chk_platform_roles_role CHECK (
        role IN (
            'super_admin',
            'admin',
            'support'
        )
    ),

    metadata JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_platform_roles_role ON platform_roles(role);

-- +goose Down
DROP INDEX IF EXISTS idx_platform_roles_role;
DROP TABLE IF EXISTS platform_roles;
