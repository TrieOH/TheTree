-- +goose Up
ALTER TABLE editions ADD COLUMN trie_payments_provider_public_key VARCHAR(256) NULL;

-- +goose Down
ALTER TABLE editions DROP COLUMN trie_payments_provider_public_key;