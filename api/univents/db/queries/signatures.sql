-- name: CreateSignature :one
INSERT INTO signatures (edition_id, created_by, signatory_name, signatory_title, signatory_email, signatory_user_id, image_url)
VALUES (@edition_id, @created_by, @signatory_name, @signatory_title, @signatory_email, @signatory_user_id, @image_url)
RETURNING *;

-- name: GetSignatureByID :one
SELECT *
FROM signatures
WHERE id = @id
  AND deleted_at IS NULL;

-- name: ListSignaturesByEdition :many
SELECT *
FROM signatures
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: DeleteSignature :exec
UPDATE signatures
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: CreateSignatureRequest :one
INSERT INTO signature_requests (edition_id, created_by, signatory_name, signatory_title, signatory_email, signatory_user_id, idempotency_key, status, expires_at)
VALUES (@edition_id, @created_by, @signatory_name, @signatory_title, @signatory_email, @signatory_user_id, @idempotency_key, 'pending', @expires_at)
RETURNING *;

-- name: GetSignatureRequestByIdempotencyKey :one
SELECT *
FROM signature_requests
WHERE idempotency_key = @idempotency_key
  AND deleted_at IS NULL;

-- name: GetSignatureRequestByID :one
SELECT *
FROM signature_requests
WHERE id = @id
  AND deleted_at IS NULL;

-- name: CompleteSignatureRequest :execrows
UPDATE signature_requests
SET
    status = 'completed',
    signature_id = @signature_id,
    updated_at = now()
WHERE id = @id
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: CancelSignatureRequest :execrows
UPDATE signature_requests
SET
    status = 'cancelled',
    status_reason = @status_reason,
    updated_at = now()
WHERE id = @id
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: ListSignatureRequestsByEdition :many
SELECT *
FROM signature_requests
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: ExpireStaleSignatureRequests :exec
UPDATE signature_requests
SET
    status = 'expired',
    updated_at = now()
WHERE status = 'pending'
  AND expires_at < now()
  AND deleted_at IS NULL;