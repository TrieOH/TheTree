-- name: CreateActivity :one
INSERT INTO activities (id, edition_id, title, description, location, starts_at, ends_at, presenter_name, token_cost, has_capacity, capacity, remaining_capacity, difficulty, created_by, scope_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: PublishActivity :exec
UPDATE activities
SET
    status = 'published'
WHERE id = $1 and status = 'draft';

-- name: StartActivity :exec
UPDATE activities
SET status = 'ongoing'
WHERE id = $1 AND status = 'published';

-- name: FinishActivity :exec
UPDATE activities
SET status = 'completed'
WHERE id = $1 AND status = 'ongoing';

-- name: GetActivityByID :one
SELECT *
FROM activities
WHERE id = $1;

-- name: ListEditionActivities :many
SELECT *
FROM activities
WHERE edition_id = $1 AND status != 'draft';

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