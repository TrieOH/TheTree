-- name: AddSignatureToEdition :one
INSERT INTO signatures (id, edition_id, title, url)
VALUES (@id, @edition_id, @title, @url)
RETURNING *;

-- name: RemoveSignatureFromEdition :exec
DELETE FROM signatures
WHERE id = @id
  AND edition_id = @edition_id;

-- name: GetSignatureByID :one
SELECT *
FROM signatures
WHERE edition_id = @edition_id
  AND id = @id;

-- name: ListSignaturesFromEdition :many
SELECT *
FROM signatures
WHERE edition_id = @edition_id;