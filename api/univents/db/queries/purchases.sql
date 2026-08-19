-- name: CreatePurchase :one
INSERT INTO purchases (edition_id, purchaser_id, status, status_reason, total_cents, currency, payment_method, payssage_seller_id, payssage_intent_id, payer_email, qr_code, qr_code_base64, expires_at, river_job_id)
VALUES (@edition_id, @purchaser_id, @status, @status_reason, @total_cents, @currency, @payment_method, @payssage_seller_id, @payssage_intent_id, @payer_email, @qr_code, @qr_code_base64, @expires_at, @river_job_id)
RETURNING *;

-- name: CreatePurchaseItem :one
INSERT INTO purchase_items (purchase_id, item_type, item_id, quantity, unit_price_cents, registration_id, product_purchase_id, participation_id)
VALUES (@purchase_id, @item_type, @item_id, @quantity, @unit_price_cents, @registration_id, @product_purchase_id, @participation_id)
RETURNING *;

-- name: GetPurchaseByID :one
SELECT *
FROM purchases
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetPurchaseByIDForOwner :one
SELECT *
FROM purchases
WHERE id = @id
  AND purchaser_id = @purchaser_id
  AND deleted_at IS NULL;

-- name: GetPurchaseByIntentID :one
SELECT *
FROM purchases
WHERE payssage_intent_id = @intent_id
  AND deleted_at IS NULL;

-- name: ListPurchasesByPurchaser :many
SELECT *
FROM purchases
WHERE purchaser_id = @purchaser_id
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListPurchasesByEdition :many
SELECT *
FROM purchases
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdatePurchaseStatus :one
UPDATE purchases
SET
    status        = @status,
    status_reason = @status_reason,
    updated_at    = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdatePurchaseStatusIf :one
-- Guarded status transition (webhook receiver, split 4): only flips when
-- the purchase is in the expected @from_status state, so a duplicate
-- delivery is a no-op (idempotent webhook). Returns no rows when the guard
-- misses.
UPDATE purchases
SET
    status        = @to_status,
    status_reason = @status_reason,
    updated_at    = now()
WHERE id = @id
  AND status = @from_status
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdatePurchaseRiverJob :one
-- Stores the expiry job's river id on the purchase (checkout, split 7): the
-- job is enqueued with river.InsertTx inside the checkout tx, then this
-- write links it so the webhook receiver can cancel it on approve.
UPDATE purchases
SET
    river_job_id = @river_job_id,
    updated_at   = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: AttachIntentToPurchase :one
-- Stores the Payssage intent on the purchase after the intent was created
-- (checkout, split 7, post-commit Tx 2): seller, intent id (the D2
-- correlation key), and the pix QR. Null for free orders.
UPDATE purchases
SET
    payssage_seller_id = @payssage_seller_id,
    payssage_intent_id = @payssage_intent_id,
    qr_code            = @qr_code,
    qr_code_base64     = @qr_code_base64,
    updated_at         = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: ListPurchaseItemsByPurchase :many
SELECT *
FROM purchase_items
WHERE purchase_id = @purchase_id;

-- Availability: available = base - reserved, where "reserved" is the
-- quantity of purchase_items joined to purchases with status IN
-- ('pending','approved'). base NULL = unlimited (never sold out). Kept as
-- three per-type queries (ticket_types / product_variants /
-- program_occurrences) so split 7 can run them inside the checkout tx
-- under the item row locks; item ids never cross tables.

-- name: ListTicketTypeAvailability :many
WITH reserved AS (
    SELECT pi.item_id, SUM(pi.quantity) AS reserved_quantity
    FROM purchase_items pi
    JOIN purchases p ON p.id = pi.purchase_id
    WHERE pi.item_type = 'ticket'
      AND p.status IN ('pending', 'approved')
      AND p.deleted_at IS NULL
    GROUP BY pi.item_id
)
SELECT
    tt.id               AS item_id,
    tt.max_quantity     AS base_quantity,
    COALESCE(r.reserved_quantity, 0)::BIGINT AS reserved_quantity
FROM ticket_types tt
LEFT JOIN reserved r ON r.item_id = tt.id
WHERE tt.edition_id = @edition_id
  AND tt.deleted_at IS NULL;

-- name: ListProductVariantAvailability :many
WITH reserved AS (
    SELECT pi.item_id, SUM(pi.quantity) AS reserved_quantity
    FROM purchase_items pi
    JOIN purchases p ON p.id = pi.purchase_id
    WHERE pi.item_type = 'product'
      AND p.status IN ('pending', 'approved')
      AND p.deleted_at IS NULL
    GROUP BY pi.item_id
)
SELECT
    pv.id               AS item_id,
    pv.stock            AS base_quantity,
    COALESCE(r.reserved_quantity, 0)::BIGINT AS reserved_quantity
FROM product_variants pv
LEFT JOIN reserved r ON r.item_id = pv.id
WHERE pv.edition_id = @edition_id
  AND pv.deleted_at IS NULL;

-- name: ListProgramOccurrenceAvailability :many
WITH reserved AS (
    SELECT pi.item_id, SUM(pi.quantity) AS reserved_quantity
    FROM purchase_items pi
    JOIN purchases p ON p.id = pi.purchase_id
    WHERE pi.item_type = 'program_occurrence'
      AND p.status IN ('pending', 'approved')
      AND p.deleted_at IS NULL
    GROUP BY pi.item_id
)
SELECT
    po.id               AS item_id,
    po.max_capacity     AS base_quantity,
    COALESCE(r.reserved_quantity, 0)::BIGINT AS reserved_quantity
FROM program_occurrences po
LEFT JOIN reserved r ON r.item_id = po.id
WHERE po.edition_id = @edition_id
  AND po.deleted_at IS NULL;
