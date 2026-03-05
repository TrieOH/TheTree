-- name: CreateTicket :one
INSERT INTO tickets (id, edition_id, name, description, created_by, scope_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: AddTicketPermission :one
INSERT INTO ticket_permissions (id, ticket_id, permission_type, activity_id, product_id, checkpoint_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: RemoveTicketPermission :exec
DELETE FROM ticket_permissions
WHERE id = $1 AND ticket_id = $2;