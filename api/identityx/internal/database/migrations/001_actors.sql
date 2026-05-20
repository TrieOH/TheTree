-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE actors (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID, -- NULL = platform-level IDX account, set = project-scoped account

    auth_method TEXT NOT NULL DEFAULT 'password',
    CONSTRAINT chk_actors_auth_method CHECK (
        auth_method IN ('password', 'google', 'github')
    ),

    email TEXT NOT NULL,
    CONSTRAINT uniq_email_per_scope_per_method UNIQUE (LOWER(email), project_id, auth_method),

    password_hash TEXT,

    type TEXT NOT NULL,
    CONSTRAINT chk_actors_type CHECK (type IN ('human', 'service', 'machine')),

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_actors_type ON actors (type);
CREATE INDEX idx_actors_created_at ON actors (created_at);
CREATE INDEX idx_actors_metadata_gin ON actors USING GIN (metadata);

CREATE TABLE actor_profiles (
    actor_id UUID PRIMARY KEY REFERENCES actors(id),

    profile JSONB NOT NULL DEFAULT '{}'::jsonb,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_actor_profiles_profile_gin ON actor_profiles USING GIN (profile);

CREATE TABLE actor_external_identities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    actor_id UUID NOT NULL REFERENCES actors(id)
        ON DELETE CASCADE,

    provider TEXT NOT NULL,
    CONSTRAINT chk_actor_external_identities_provider CHECK (
        provider IN ('google', 'github')
    ),

    -- stable ID from the provider, never changes even if email does
    subject TEXT NOT NULL,
    CONSTRAINT uniq_external_identity UNIQUE (provider, subject),

    email TEXT,

    encrypted_access_token TEXT,
    encrypted_refresh_token TEXT,
    token_expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_actor_external_identities_actor_id
    ON actor_external_identities (actor_id);

-- +goose Down
DROP INDEX IF EXISTS idx_actor_external_identities_actor_id;
DROP TABLE IF EXISTS actor_external_identities;
DROP INDEX IF EXISTS idx_actor_profiles_profile_gin;
DROP TABLE IF EXISTS actor_profiles;
DROP INDEX IF EXISTS idx_actors_metadata_gin;
DROP INDEX IF EXISTS idx_actors_created_at;
DROP INDEX IF EXISTS idx_actors_type;
DROP TABLE IF EXISTS actors;
DROP EXTENSION IF EXISTS "pgcrypto";