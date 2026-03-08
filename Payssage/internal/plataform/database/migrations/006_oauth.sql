-- +goose Up
CREATE TABLE oauth_states (
    state             TEXT PRIMARY KEY,
    workspace_id      UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    provider          TEXT NOT NULL,
    flow              TEXT NOT NULL, -- 'setup' or 'connect'
    is_marketplace    BOOL NOT NULL DEFAULT false,
    fee_bps           INT NOT NULL DEFAULT 0,
    final_redirect_url TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at        TIMESTAMPTZ NOT NULL
);

CREATE TABLE provider_credentials (
    id           UUID PRIMARY KEY DEFAULT uuidv7(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    provider     TEXT NOT NULL,
    display_name TEXT NOT NULL,
    credentials  JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at   TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX idx_provider_credentials_active
ON provider_credentials (workspace_id, provider, display_name)
WHERE revoked_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_provider_credentials_active;
DROP TABLE IF EXISTS provider_credentials;
DROP TABLE IF EXISTS oauth_states;