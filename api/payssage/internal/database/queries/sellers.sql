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

-- name: ListSellersByWallet :many
SELECT *
FROM sellers
WHERE wallet_id = @wallet_id;

-- name: RevokeSeller :exec
UPDATE sellers
SET revoked_at = NOW(),
    credentials = '{}'
WHERE id = @id
  AND revoked_at IS NULL;
