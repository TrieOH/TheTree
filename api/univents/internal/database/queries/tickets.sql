-- name: CreateTicket :one
INSERT INTO tickets (id, edition_id, name, description, created_by, scope_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTicketByID :one
SELECT *
FROM tickets
WHERE id = $1;

-- name: AddTicketPermission :one
INSERT INTO ticket_permissions (id, ticket_id, permission_type, activity_id, product_id, checkpoint_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: RemoveTicketPermission :exec
DELETE FROM ticket_permissions
WHERE id = $1 AND ticket_id = $2;

-- name: ListEditionTickets :many
SELECT *
FROM tickets
WHERE edition_id = $1;

-- name: GetTicketGrantsByPaymentIntent :many
SELECT
    purchase_items.item_id as ticket_id,
    COALESCE(purchase_items.assigned_to_user_id, purchases.user_id) AS user_id
FROM purchase_items
         JOIN purchases ON purchases.id = purchase_items.purchase_id
WHERE purchases.payment_id = sqlc.arg(payment_id)
  AND purchase_items.item_type = 'ticket';

-- name: GetTicketPermissions :many
SELECT *
FROM ticket_permissions
WHERE ticket_id = $1;