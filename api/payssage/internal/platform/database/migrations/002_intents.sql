-- +goose Up
CREATE TYPE intent_status AS ENUM (
    'pending',
    'succeeded',
    'cancelled',
    'failed'
);

CREATE TABLE intents (
    id            UUID PRIMARY KEY DEFAULT uuidv7(),
    workspace_id  UUID NOT NULL,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    currency      CHAR(3) NOT NULL,
    status        intent_status NOT NULL DEFAULT 'pending',
    provider      TEXT NOT NULL,
    provider_data JSONB NOT NULL DEFAULT '{}'::JSONB,
    metadata      JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_intents_workspace_id ON intents (workspace_id);
CREATE INDEX idx_intents_provider_data ON intents USING GIN (provider_data);

-- +goose Down
DROP INDEX IF EXISTS idx_intents_provider_data;
DROP INDEX IF EXISTS idx_intents_workspace_id;
DROP TABLE IF EXISTS intents;
DROP TYPE IF EXISTS intent_status;