-- name: CreateProduct :one
INSERT INTO products (id, scope_id, edition_id, name, description, type, price_cents, available_from, available_until, has_inventory, inventory_quantity, inventory_remaining, created_by, ticket_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;