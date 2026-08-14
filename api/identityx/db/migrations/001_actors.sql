-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE actors(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID, -- NULL = platform-level IDX account, set = project-scoped account

    auth_method TEXT NOT NULL DEFAULT 'password',
    CONSTRAINT chk_actors_auth_method CHECK (
        auth_method IN ('api_key', 'password', 'google', 'github')
    ),

    verified_at TIMESTAMPTZ,
    password_hash TEXT,
    email TEXT,
    CONSTRAINT chk_actors_email_required_for_humans CHECK (type != 'human' OR email IS NOT NULL),

    type TEXT NOT NULL,
    CONSTRAINT chk_actors_type CHECK (type IN ('human', 'service', 'machine')),

    metadata JSONB DEFAULT '{}'::jsonb,

    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX uniq_email_per_scope_per_method ON actors (LOWER(email), project_id, auth_method) NULLS NOT DISTINCT;

CREATE INDEX idx_actors_type ON actors (type);
CREATE INDEX idx_actors_created_at ON actors (created_at);
CREATE INDEX idx_actors_metadata_gin ON actors USING GIN (metadata);

CREATE TABLE actor_profiles(
    actor_id UUID PRIMARY KEY REFERENCES actors(id),

    profile JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- schema version the profile was validated against; outdated flags
    -- documents that failed to migrate to the active schema
    schema_version INTEGER NOT NULL DEFAULT 1,
    outdated BOOLEAN NOT NULL DEFAULT false,

    handle TEXT, -- unique when present: NULLs (no handle set) never collide
    pfp_url TEXT, -- profile picture URL, first-class column

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uniq_actor_profiles_handle
    ON actor_profiles (handle) WHERE handle IS NOT NULL;

CREATE INDEX idx_actor_profiles_profile_gin ON actor_profiles USING GIN (profile);

CREATE TABLE actor_external_identities(
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
    ON actor_external_identities(actor_id);

-- +goose Down
DROP INDEX IF EXISTS idx_actor_external_identities_actor_id;
DROP TABLE IF EXISTS actor_external_identities;
DROP INDEX IF EXISTS idx_actor_profiles_profile_gin;
DROP INDEX IF EXISTS uniq_actor_profiles_handle;
DROP TABLE IF EXISTS actor_profiles;
DROP INDEX IF EXISTS idx_actors_metadata_gin;
DROP INDEX IF EXISTS idx_actors_created_at;
DROP INDEX IF EXISTS idx_actors_type;
DROP TABLE IF EXISTS actors;
DROP EXTENSION IF EXISTS "pgcrypto";