-- +goose Up
CREATE SEQUENCE intents_seq;

CREATE TABLE intents (
    id        TEXT PRIMARY KEY DEFAULT 'pi_' || nextval('intents_seq'),
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT,
    seller_id UUID REFERENCES sellers(id) ON DELETE RESTRICT,

    amount_cents        BIGINT NOT NULL,
    currency      CHAR(3) NOT NULL,
    status        TEXT NOT NULL DEFAULT 'pending',
    provider_data JSONB NOT NULL DEFAULT '{}'::JSONB,
    metadata      JSONB NOT NULL DEFAULT '{}'::JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_intents_amount_cents CHECK (amount_cents > 0),
    CONSTRAINT chk_intents_status CHECK (status IN ('pending', 'processing', 'succeeded', 'cancelled', 'failed'))
);

CREATE INDEX idx_intents_wallet_id ON intents (wallet_id);
CREATE INDEX idx_intents_seller_id ON intents (seller_id);
CREATE INDEX idx_intents_status ON intents (status);

-- +goose Down
DROP INDEX IF EXISTS idx_intents_status;
DROP INDEX IF EXISTS idx_intents_seller_id;
DROP INDEX IF EXISTS idx_intents_wallet_id;
DROP TABLE IF EXISTS intents;
DROP SEQUENCE IF EXISTS intents_seq;