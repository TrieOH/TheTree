-- +goose Up
ALTER TABLE events ADD COLUMN payssage_public_key TEXT;

-- The wallet is the event's permanent payment container: it may exist
-- without a seller (connected account). Relax the both-or-neither rule to a
-- one-way rule (seller requires wallet, wallet alone is fine).
ALTER TABLE events DROP CONSTRAINT chk_event_payments_config_complete;

ALTER TABLE events
    ADD CONSTRAINT chk_event_payments_config_complete CHECK (
        payssage_seller_id IS NULL OR payssage_wallet_id IS NOT NULL
    );
-- +goose Down
ALTER TABLE events DROP CONSTRAINT chk_event_payments_config_complete;

ALTER TABLE events
    ADD CONSTRAINT chk_event_payments_config_complete CHECK (
        (payssage_seller_id IS NULL) = (payssage_wallet_id IS NULL)
    );

ALTER TABLE events DROP COLUMN payssage_public_key;
