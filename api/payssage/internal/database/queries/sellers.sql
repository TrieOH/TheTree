-- name: CreateSeller :one
INSERT INTO sellers (wallet_id, provider, provider_user_id, credentials)
VALUES (@wallet_id, @provider, @provider_user_id, @credentials)
RETURNING *;

-- name: GetSellerByID :one
SELECT *
FROM sellers
WHERE id = @id;

-- name: ListSellers :many
SELECT *
FROM sellers;
