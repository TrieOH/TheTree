-- +goose Up
CREATE TABLE marketplace_configs (
    id            UUID PRIMARY KEY DEFAULT uuidv7(),
    workspace_id  UUID NOT NULL UNIQUE REFERENCES workspaces(id) ON DELETE CASCADE,
    credential_id UUID NOT NULL REFERENCES provider_credentials(id),
    fee_bps       INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS marketplace_configs;