-- name: WaitlistAppend :one
INSERT INTO waitlist_entries (activity_id, user_id, position)
SELECT $1, $2, COALESCE(MAX(position) + 1, 0)
FROM waitlist_entries
WHERE activity_id = $1
RETURNING *;

-- name: WaitlistGetFirst :one
SELECT * FROM waitlist_entries
WHERE activity_id = $1
ORDER BY position
LIMIT 1;

-- name: WaitlistRemove :exec
DELETE FROM waitlist_entries WHERE id = $1;

-- name: WaitlistGetOrdered :many
SELECT
    w.*,
    ROW_NUMBER() OVER (ORDER BY w.position) as queue_position
FROM waitlist_entries w
WHERE activity_id = $1
ORDER BY w.position;