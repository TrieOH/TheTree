-- +goose Up
CREATE TABLE webhook_endpoints (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    scope_id UUID NOT NULL,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_webhook_endpoints_workspace_id ON webhook_endpoints (workspace_id);

CREATE TYPE delivery_status AS ENUM (
    'pending',
    'delivered',
    'failed'
);

CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    endpoint_id UUID NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    intent_id UUID NOT NULL REFERENCES intents(id) ON DELETE CASCADE,
    event TEXT NOT NULL,
    payload JSONB NOT NULL,
    status delivery_status NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    last_attempted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE webhook_events (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    workspace_id UUID NULL REFERENCES workspaces(id) ON DELETE SET NULL,
    intent_id    UUID NULL REFERENCES intents(id) ON DELETE SET NULL,
    provider     TEXT NOT NULL,
    external_id  TEXT NULL,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_webhook_deliveries_endpoint_id ON webhook_deliveries (endpoint_id);
CREATE INDEX idx_webhook_deliveries_intent_id ON webhook_deliveries (intent_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries (status);
-- +goose Down
DROP INDEX IF EXISTS idx_webhook_deliveries_status;
DROP INDEX IF EXISTS idx_webhook_deliveries_intent_id;
DROP INDEX IF EXISTS idx_webhook_deliveries_endpoint_id;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TYPE IF EXISTS delivery_status;
DROP INDEX IF EXISTS idx_webhook_endpoints_workspace_id;
DROP TABLE IF EXISTS webhook_endpoints;
