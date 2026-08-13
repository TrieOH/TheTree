-- 011_purchase_status_failed_rejected.sql
-- +goose Up
-- Distinguish *why* a purchase died (D4): the customer cancelling is
-- "cancelled", while a declined payment is "failed" (payment.failed) or
-- "rejected" (payment.rejected — e.g. MP's risk engine). Previously all
-- three collapsed into "cancelled", which mislead the buyer.
ALTER TABLE purchases DROP CONSTRAINT chk_purchases_status_valid;
ALTER TABLE purchases ADD CONSTRAINT chk_purchases_status_valid CHECK (
    status IN ('pending', 'approved', 'expired', 'cancelled', 'failed', 'rejected')
);

-- +goose Down
ALTER TABLE purchases DROP CONSTRAINT chk_purchases_status_valid;
ALTER TABLE purchases ADD CONSTRAINT chk_purchases_status_valid CHECK (
    status IN ('pending', 'approved', 'expired', 'cancelled')
);
