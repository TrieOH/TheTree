-- +goose Up
CREATE TABLE encryption_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    project_id UUID REFERENCES projects(id)
        ON DELETE SET NULL,

    type TEXT NOT NULL,
    CONSTRAINT chk_encryption_keys_type CHECK (
        type IN (
            'jwt',
            'encryption',
            'signing'
        )
    ),

    public_key TEXT NOT NULL,
    encrypted_private_key TEXT NOT NULL,

    algorithm TEXT NOT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rotated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_encryption_keys_project_id ON encryption_keys (project_id);
CREATE INDEX idx_encryption_keys_active ON encryption_keys (active);
CREATE INDEX idx_encryption_keys_type ON encryption_keys (type);

-- +goose Down
DROP INDEX IF EXISTS idx_encryption_keys_type;
DROP INDEX IF EXISTS idx_encryption_keys_active;
DROP INDEX IF EXISTS idx_encryption_keys_project_id;
DROP TABLE IF EXISTS encryption_keys;