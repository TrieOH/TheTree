-- +goose Up

CREATE TYPE identity_type AS ENUM ('client', 'project');

CREATE TABLE session_identities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    type identity_type NOT NULL,
    entity_id UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (type, entity_id)
);

ALTER TABLE sessions
    ADD COLUMN identity_id UUID NOT NULL
        REFERENCES session_identities(id)
            ON DELETE CASCADE;

CREATE INDEX idx_sessions_identity_id
    ON sessions(identity_id);

-- +goose Down

DROP INDEX IF EXISTS idx_sessions_identity_id;

ALTER TABLE sessions
DROP COLUMN IF EXISTS identity_id;

DROP TABLE IF EXISTS session_identities;

DROP TYPE IF EXISTS identity_type;