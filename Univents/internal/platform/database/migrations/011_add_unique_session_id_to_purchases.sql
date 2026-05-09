-- +goose Up
ALTER TABLE purchases ADD CONSTRAINT uq_purchases_session_id UNIQUE (session_id);

-- +goose Down
ALTER TABLE purchases DROP CONSTRAINT IF EXISTS uq_purchases_session_id;