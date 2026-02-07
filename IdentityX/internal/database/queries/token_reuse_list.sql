-- name: TokenReuseListAppend :exec
INSERT INTO token_reuse_list (jit, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: TokenReuseListExists :one
SELECT EXISTS (
    SELECT 1
    FROM token_reuse_list
    WHERE jit = $1
      AND user_id = $2
);

-- name: DeleteExpiredTokenReuseListEntries :exec
DELETE FROM token_reuse_list
WHERE expired_at < now();