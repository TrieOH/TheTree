-- +goose Up
-- Refund flow (refund plan slices 3-4):
--   1. purchases gain the `refunded` status (approved -> refunded is
--      webhook-confirmed via payment.refunded).
--   2. purchases gain payer_email: the person the payment provider will
--      refund (set at checkout from the checkout payer), used by the
--      organizer orders read to identify who gets the money.
ALTER TABLE purchases DROP CONSTRAINT chk_purchases_status_valid;
ALTER TABLE purchases ADD CONSTRAINT chk_purchases_status_valid CHECK (
    status IN ('pending', 'approved', 'expired', 'cancelled', 'failed', 'rejected', 'refunded')
);

ALTER TABLE purchases ADD COLUMN payer_email TEXT;

-- +goose Down
ALTER TABLE purchases DROP COLUMN payer_email;
ALTER TABLE purchases DROP CONSTRAINT chk_purchases_status_valid;
ALTER TABLE purchases ADD CONSTRAINT chk_purchases_status_valid CHECK (
    status IN ('pending', 'approved', 'expired', 'cancelled', 'failed', 'rejected')
);
