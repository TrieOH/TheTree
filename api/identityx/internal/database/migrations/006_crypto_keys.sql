-- +goose Up
CREATE TABLE crypto_keys (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    project_id UUID REFERENCES projects(id)
        ON DELETE CASCADE,

    type TEXT NOT NULL,
    CONSTRAINT chk_crypto_keys_type CHECK (
        type IN (
            'encryption',
            'signing'
        )
    ),

    status TEXT NOT NULL DEFAULT 'active',
    CONSTRAINT chk_crypto_keys_status CHECK (
        status IN (
           'active',    -- current, used for new signing/encrypting
           'retiring',  -- rotation in progress, still verifies/decrypts old tokens
           'retired',   -- gracefully phased out, no longer verifies (tokens expired naturally)
           'revoked'    -- forcefully killed, reject everything signed with this
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

CREATE INDEX idx_crypto_keys_project_id ON crypto_keys (project_id);
CREATE INDEX idx_crypto_keys_active ON crypto_keys (active);
CREATE INDEX idx_crypto_keys_type ON crypto_keys (type);

-- +goose Down
DROP INDEX IF EXISTS idx_crypto_keys_type;
DROP INDEX IF EXISTS idx_crypto_keys_active;
DROP INDEX IF EXISTS idx_crypto_keys_project_id;
DROP TABLE IF EXISTS crypto_keys;