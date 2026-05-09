-- name: CreateCheckpoint :one
INSERT INTO checkpoints (id, scope_id, edition_id, name, type, access_mode, created_by, starts_at, ends_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListEditionCheckpoints :many
SELECT *
FROM checkpoints
WHERE edition_id = $1;

-- name: GetCheckpointByID :one
SELECT *
FROM checkpoints
WHERE id = $1;