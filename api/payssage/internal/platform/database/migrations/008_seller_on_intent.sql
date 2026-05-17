-- +goose Up
ALTER TABLE intents ADD COLUMN seller_credential_id UUID;

CREATE INDEX idx_intents_seller_credential_id ON intents (seller_credential_id);

-- +goose Down
DROP INDEX IF EXISTS idx_intents_seller_credential_id;
ALTER TABLE intents DROP COLUMN seller_credential_id;