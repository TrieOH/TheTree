-- name: CreateTicketType :one
INSERT INTO ticket_types (edition_id, name, description, access_level, price, max_quantity)
VALUES (@edition_id, @name, @description, @access_level, @price, @max_quantity)
RETURNING *;

-- name: GetTicketTypeByID :one
SELECT *
FROM ticket_types
WHERE id = @id;

-- name: ListTicketTypesByEdition :many
SELECT *
FROM ticket_types
WHERE edition_id = @edition_id
ORDER BY created_at;

-- name: PatchTicketType :one
UPDATE ticket_types
SET
    name         = @name,
    description  = @description,
    access_level = @access_level,
    price        = @price,
    max_quantity = @max_quantity,
    updated_at   = now()
WHERE id = @id
RETURNING *;
