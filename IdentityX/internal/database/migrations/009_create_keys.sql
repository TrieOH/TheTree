-- +goose Up

CREATE TYPE key_type_enum AS ENUM (
    'goauth',
    'project'
);

CREATE TYPE key_usage_enum AS ENUM (
    'sign',
    'verify'
);

CREATE TYPE key_status_enum AS ENUM (
    'active',
    'rotated',
    'revoked'
);

CREATE TABLE key_pair (
    id UUID PRIMARY KEY DEFAULT uuidv7(),

    kid TEXT NOT NULL UNIQUE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,

    key_type key_type_enum NOT NULL,
    algorithm TEXT NOT NULL DEFAULT 'Ed25519',

    public_key TEXT NOT NULL,
    private_key BYTEA NOT NULL, -- envelope-encrypted, never plaintext

    usage key_usage_enum NOT NULL DEFAULT 'sign',
    status key_status_enum NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,

    verify_expires_at TIMESTAMPTZ NOT NULL,

    CHECK (
        (key_type = 'goauth' AND project_id IS NULL)
            OR
        (key_type = 'project' AND project_id IS NOT NULL)
    ),

    CHECK (
        NOT (usage = 'sign' AND status = 'rotated')
    )
);

CREATE INDEX idx_key_pair_project_active_sign
    ON key_pair (project_id, created_at DESC)
    WHERE status = 'active' AND usage = 'sign';

CREATE INDEX idx_key_pair_kid_lookup
    ON key_pair (kid);

CREATE INDEX idx_key_pair_project_jwks
    ON key_pair (project_id)
    WHERE status IN ('active', 'rotated');

CREATE INDEX idx_key_pair_goauth_active_sign
    ON key_pair (created_at DESC)
    WHERE key_type = 'goauth'
      AND project_id IS NULL
      AND status = 'active'
      AND usage = 'sign';

CREATE INDEX idx_key_pair_goauth_jwks
    ON key_pair (created_at DESC)
    WHERE key_type = 'goauth'
      AND project_id IS NULL
      AND status IN ('active', 'rotated');

CREATE UNIQUE INDEX uniq_goauth_single_active_signing_key
    ON key_pair (key_type)
    WHERE
    key_type = 'goauth'
  AND project_id IS NULL
  AND usage = 'sign'
  AND status = 'active';

CREATE INDEX idx_key_pair_project_jwks
    ON key_pair (project_id)
    WHERE status IN ('active', 'rotated')
      AND verify_expires_at > now();

-- +goose Down
DROP INDEX IF EXISTS idx_key_pair_project_jwks;
DROP INDEX IF EXISTS uniq_goauth_single_active_signing_key;
DROP INDEX IF EXISTS idx_key_pair_goauth_jwks;
DROP INDEX IF EXISTS idx_key_pair_goauth_active_sign;
DROP INDEX IF EXISTS idx_key_pair_project_jwks;
DROP INDEX IF EXISTS idx_key_pair_kid_lookup;
DROP INDEX IF EXISTS idx_key_pair_project_active_sign;
DROP TABLE IF EXISTS key_pair;
