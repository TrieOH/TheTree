-- name: CreateTicket :one
INSERT INTO tickets (id, edition_id, name, description, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
