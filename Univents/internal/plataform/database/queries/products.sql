-- name: CreateProduct :one
INSERT INTO products (id, scope_id, edition_id, name, description, type, price_cents, available_from, available_until, has_inventory, inventory_quantity, inventory_remaining, created_by, ticket_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: PublishProduct :exec
UPDATE products
SET
    status = 'available'
WHERE id = $1 and status = 'draft';

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: GetProductsByIDs :many
SELECT *
FROM products
WHERE products.id = ANY(sqlc.arg(ids)::uuid[])
AND products.deleted_at IS NULL;

-- name: ListEditionProducts :many
SELECT *
FROM products
WHERE edition_id = $1 AND status != 'draft';

-- name: ListEditionProductsAdmin :many
SELECT *
FROM products
WHERE edition_id = $1;

-- name: ReserveProductNoInventory :exec
INSERT INTO product_reservations (session_id, product_id, quantity, expires_at)
VALUES (sqlc.arg(session_id), sqlc.arg(product_id), sqlc.arg(quantity), sqlc.arg(expires_at));

-- name: ReserveProduct :one
WITH available AS (
    SELECT inventory_remaining
    FROM products
    WHERE id = sqlc.arg(product_id)
      AND has_inventory = TRUE
      AND inventory_remaining > 0
    FOR UPDATE
    ),
updated AS (
    UPDATE products
    SET inventory_remaining = inventory_remaining - LEAST(sqlc.arg(quantity)::int, (SELECT inventory_remaining FROM available))
    WHERE id = sqlc.arg(product_id)
      AND EXISTS (SELECT 1 FROM available)
    RETURNING id, inventory_remaining
)
INSERT INTO product_reservations (session_id, product_id, quantity, expires_at)
SELECT
    sqlc.arg(session_id),
    sqlc.arg(product_id),
    LEAST(sqlc.arg(quantity)::int, available.inventory_remaining),
    sqlc.arg(expires_at)
FROM available
WHERE EXISTS (SELECT 1 FROM updated)
    RETURNING
    (SELECT inventory_remaining FROM updated) AS inventory_remaining,
    quantity AS reserved_quantity;

-- name: UnreserveProducts :many
WITH deleted AS (
DELETE FROM product_reservations
WHERE product_reservations.session_id = sqlc.arg(session_id)
    RETURNING product_reservations.product_id, product_reservations.quantity
)
UPDATE products
SET inventory_remaining = inventory_remaining + deleted.quantity
    FROM deleted
WHERE products.id = deleted.product_id
  AND products.has_inventory = TRUE
    RETURNING products.id, products.inventory_remaining;

-- name: DeleteReservation :exec
DELETE FROM product_reservations
WHERE session_id = $1;

-- name: ConfirmPurchase :exec
UPDATE purchases
SET
    status = 'completed',
    fulfilled_at = now()
WHERE payment_id = $1;

-- name: CancelPurchase :exec
UPDATE purchases
SET status = 'cancelled'
WHERE payment_id = $1 AND status = 'pending';

-- name: GetReservationItems :many
SELECT
    product_reservations.session_id,
    product_reservations.product_id,
    product_reservations.quantity,
    products.price_cents,
    products.type,
    products.ticket_id
FROM product_reservations
JOIN products ON products.id = product_reservations.product_id
WHERE product_reservations.session_id = sqlc.arg(session_id);

-- name: CreatePurchase :one
INSERT INTO purchases (edition_id, session_id, user_id, status, subtotal_cents, total_cents, payment_provider, payment_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (session_id) DO NOTHING
RETURNING *;

-- name: CreatePurchaseItem :one
INSERT INTO purchase_items (purchase_id, item_type, item_id, quantity, unit_price_cents, total_price_cents)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPurchaseByPaymentID :one
SELECT * FROM purchases
WHERE payment_id = $1
  AND deleted_at IS NULL;

-- name: GetPurchaseBySessionID :one
SELECT * FROM purchases
WHERE session_id = $1
  AND deleted_at IS NULL;

-- name: ListUserPurchases :many
SELECT * FROM purchases
WHERE user_id = $1;

-- name: ListPurchaseItems :many
SELECT pi.*
FROM purchase_items pi
JOIN purchases p ON p.id = pi.purchase_id
WHERE pi.purchase_id = $1
  AND p.user_id = $2;