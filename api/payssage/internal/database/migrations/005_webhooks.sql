-- 005_webhooks.sql
-- +goose Up
CREATE TABLE webhook_endpoints (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id  UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    secret     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE webhook_events (
    id          UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id   UUID REFERENCES wallets(id) ON DELETE SET NULL,
    intent_id   TEXT REFERENCES intents(id) ON DELETE SET NULL,
    provider    TEXT NOT NULL,
    external_id TEXT,
    event_type  TEXT NOT NULL,
    payload     JSONB NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE webhook_deliveries (
    id                UUID PRIMARY KEY DEFAULT uuidv7(),
    endpoint_id       UUID NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    event_id          UUID NOT NULL REFERENCES webhook_events(id) ON DELETE CASCADE,
    status            TEXT NOT NULL DEFAULT 'pending',
    attempts          INT NOT NULL DEFAULT 0,
    last_attempted_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_webhook_deliveries_status CHECK (status IN ('pending', 'delivered', 'failed'))
);

CREATE UNIQUE INDEX uniq_webhook_events_external_id
    ON webhook_events (provider, external_id)
    WHERE external_id IS NOT NULL;

CREATE INDEX idx_webhook_endpoints_wallet_id ON webhook_endpoints (wallet_id);
CREATE INDEX idx_webhook_events_intent_id ON webhook_events (intent_id);
CREATE INDEX idx_webhook_deliveries_endpoint_id ON webhook_deliveries (endpoint_id);
CREATE INDEX idx_webhook_deliveries_event_id ON webhook_deliveries (event_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries (status);

-- +goose Down
DROP INDEX IF EXISTS idx_webhook_deliveries_status;
DROP INDEX IF EXISTS idx_webhook_deliveries_event_id;
DROP INDEX IF EXISTS idx_webhook_deliveries_endpoint_id;
DROP INDEX IF EXISTS idx_webhook_events_intent_id;
DROP INDEX IF EXISTS uniq_webhook_events_external_id;
DROP INDEX IF EXISTS idx_webhook_endpoints_wallet_id;
DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS webhook_endpoints;