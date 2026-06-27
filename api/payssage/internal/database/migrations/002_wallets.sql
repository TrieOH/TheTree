-- +goose Up
CREATE TABLE wallets (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    owner_id        UUID NOT NULL,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    sandbox         BOOLEAN NOT NULL DEFAULT false,
    fee_bps         INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_wallets_fee_bps CHECK (fee_bps >= 0)
);

CREATE UNIQUE INDEX uniq_wallets_org_name
    ON wallets (organization_id, name)
    WHERE organization_id IS NOT NULL;

CREATE UNIQUE INDEX uniq_wallets_personal_name
    ON wallets (owner_id, name)
    WHERE organization_id IS NULL;

CREATE INDEX idx_wallets_owner_id ON wallets (owner_id);

-- +goose Down
DROP INDEX IF EXISTS idx_wallets_owner_id;
DROP INDEX IF EXISTS uniq_wallets_personal_name;
DROP INDEX IF EXISTS uniq_wallets_org_name;
DROP TABLE IF EXISTS wallets;