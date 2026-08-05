-- name: InsertActionToken :one
INSERT INTO action_tokens (jti, purpose, actor_id, expires_at)
VALUES (@jti, @purpose, @actor_id, @expires_at)
RETURNING *;

-- name: GetActionTokenByJTI :one
SELECT *
FROM action_tokens
WHERE jti = @jti;

-- name: ConsumeActionToken :one
UPDATE action_tokens
SET used_at = NOW()
WHERE jti = @jti
  AND used_at IS NULL
RETURNING *;

-- name: DeleteExpiredActionTokens :exec
DELETE FROM action_tokens WHERE expires_at <= NOW();
