-- +goose Up
CREATE TABLE oauth_states (
    state              TEXT PRIMARY KEY,
    wallet_id          UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    provider           TEXT NOT NULL,
    flow               TEXT NOT NULL,
    final_redirect_url TEXT NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at         TIMESTAMPTZ NOT NULL,

    CONSTRAINT chk_oauth_states_flow CHECK (flow IN ('setup', 'connect'))
);

CREATE TABLE provider_credentials (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id        UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    provider         TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,
    credentials      JSONB NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at       TIMESTAMPTZ,

    CONSTRAINT uniq_provider_credentials_active UNIQUE (wallet_id, provider)
);

CREATE TABLE sellers (
    id               UUID PRIMARY KEY DEFAULT uuidv7(),
    wallet_id        UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    provider         TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at       TIMESTAMPTZ,

    CONSTRAINT uniq_sellers_active UNIQUE (wallet_id, provider_user_id)
);

CREATE INDEX idx_provider_credentials_wallet_id ON provider_credentials (wallet_id);
CREATE INDEX idx_sellers_wallet_id ON sellers (wallet_id);

-- +goose Down
DROP INDEX IF EXISTS idx_sellers_wallet_id;
DROP INDEX IF EXISTS idx_provider_credentials_wallet_id;
DROP TABLE IF EXISTS sellers;
DROP TABLE IF EXISTS provider_credentials;
DROP TABLE IF EXISTS oauth_states;