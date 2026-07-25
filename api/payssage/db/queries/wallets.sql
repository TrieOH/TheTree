-- name: CreateWallet :one
INSERT INTO wallets (owner_id, organization_id, name, sandbox, fee_bps)
VALUES (@owner_id, @organization_id, @name, @sandbox, @fee_bps)
    RETURNING *;

-- name: GetWalletByID :one
SELECT *
FROM wallets
WHERE id = @id;

-- name: ListWallets :many
SELECT *
FROM wallets
WHERE owner_id = @owner_id
  AND organization_id IS NULL;

-- name: ListOrgWallets :many
SELECT *
FROM wallets
WHERE organization_id = @organization_id;

-- name: SetWalletSandboxState :exec
UPDATE wallets
SET
    sandbox = @sandbox
WHERE id = @id;

-- name: SetWalletFeeBPS :exec
UPDATE wallets
SET
    fee_bps = @fee_bps
WHERE id = @id;

-- name: BindCollector :exec
UPDATE wallets
SET collector_id = @collector_id
WHERE id = @id
  AND collector_id IS NULL;

-- name: UnbindCollector :exec
UPDATE wallets
SET collector_id = NULL
WHERE id = @id;