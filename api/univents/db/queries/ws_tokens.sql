-- name: CreateWsToken :one
INSERT INTO ws_tokens (purchase_id, user_id, token_hash, expires_at)
VALUES (@purchase_id, @user_id, @token_hash, @expires_at)
RETURNING *;

-- name: ConsumeWsToken :one
-- Atomic one-time consume (split 6): marks the token used only when it
-- exists, is unused, and unexpired. Returns no rows when the guard misses
-- (missing / already used / expired) — the handshake rejects.
UPDATE ws_tokens
SET used_at = now()
WHERE token_hash = @token_hash
  AND used_at IS NULL
  AND expires_at > now()
RETURNING *;
