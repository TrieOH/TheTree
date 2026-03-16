-- +goose Up
DROP INDEX IF EXISTS idx_provider_credentials_active;

CREATE UNIQUE INDEX idx_provider_credentials_active
ON provider_credentials (workspace_id, (credentials->>'provider_user_id'), display_name)
WHERE revoked_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_provider_credentials_active;

CREATE UNIQUE INDEX idx_provider_credentials_active
ON provider_credentials (workspace_id, provider, display_name)
WHERE revoked_at IS NULL;