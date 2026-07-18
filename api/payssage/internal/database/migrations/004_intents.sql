-- +goose Up
CREATE TABLE intents (
    id            UUID PRIMARY KEY,
    wallet_id     UUID NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT,
    seller_id     UUID NOT NULL REFERENCES sellers(id) ON DELETE RESTRICT,
    collector_id  UUID REFERENCES collectors(id) ON DELETE SET NULL,

    sandbox       BOOLEAN NOT NULL DEFAULT FALSE,
    amount_cents  BIGINT NOT NULL,
    currency      CHAR(3) NOT NULL,
    status        TEXT NOT NULL DEFAULT 'pending',
    provider      TEXT NOT NULL,
    provider_data JSONB NOT NULL DEFAULT '{}'::JSONB,
    metadata      JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_intents_amount_cents CHECK (amount_cents > 0),
    CONSTRAINT chk_intents_status CHECK (status IN ('pending', 'processing', 'succeeded', 'cancelled', 'failed'))
);

CREATE INDEX idx_intents_wallet_id ON intents (wallet_id);
CREATE INDEX idx_intents_seller_id ON intents (seller_id);
CREATE INDEX idx_intents_collector_id ON intents (collector_id);
CREATE INDEX idx_intents_status ON intents (status);

-- +goose Down
DROP INDEX IF EXISTS idx_intents_status;
DROP INDEX IF EXISTS idx_intents_collector_id;
DROP INDEX IF EXISTS idx_intents_seller_id;
DROP INDEX IF EXISTS idx_intents_wallet_id;
DROP TABLE IF EXISTS intents;
