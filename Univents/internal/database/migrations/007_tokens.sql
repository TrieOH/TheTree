-- +goose Up
-- User token balances
CREATE TABLE user_token_balances (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,

    balance INT NOT NULL DEFAULT 0,
    total_earned INT NOT NULL DEFAULT 0,
    total_spent INT NOT NULL DEFAULT 0,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(user_id, edition_id)
);

CREATE INDEX idx_token_balances_user ON user_token_balances(user_id);
CREATE INDEX idx_token_balances_edition ON user_token_balances(edition_id);

-- Token transactions (audit trail)
CREATE TYPE token_transaction_type AS ENUM (
    'purchased',    -- bought token pack
    'granted',      -- admin grant
    'spent',        -- used for activity/product
    'refunded',     -- activity cancelled
    'expired'       -- cleanup
);

CREATE TABLE token_transactions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    edition_id UUID NOT NULL REFERENCES editions(id) ON DELETE CASCADE,

    type token_transaction_type NOT NULL,
    amount INT NOT NULL, -- positive = credit, negative = debit

    -- source
    purchase_item_id UUID NULL REFERENCES purchase_items(id),
    activity_id UUID NULL REFERENCES activities(id),

    balance_after INT NOT NULL,

    metadata JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_token_transactions_user ON token_transactions(user_id, edition_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_token_transactions_user;
DROP TABLE IF EXISTS token_transactions;
DROP TYPE IF EXISTS token_transaction_type;
DROP INDEX IF EXISTS idx_token_balances_edition;
DROP INDEX IF EXISTS idx_token_balances_user;
DROP TABLE IF EXISTS user_token_balances;