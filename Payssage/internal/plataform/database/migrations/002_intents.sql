-- +goose Up
CREATE TYPE intent_status AS ENUM (
    'pending',
    'succeeded',
    'cancelled',
    'failed'
);

CREATE TABLE intents (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    workspace_id UUID NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    status intent_status NOT NULL DEFAULT 'pending',
    client_secret TEXT NOT NULL,
    provider TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_intents_workspace_id ON intents (workspace_id);
-- +goose Down
DROP INDEX IF EXISTS idx_intents_workspace_id;
DROP TABLE IF EXISTS intents;
DROP TYPE IF EXISTS intent_status;