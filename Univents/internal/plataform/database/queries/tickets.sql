-- name: CreateTicket :one
INSERT INTO tickets (edition_id, name, description, price_cents, has_limited_quantity, quantity_available)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
