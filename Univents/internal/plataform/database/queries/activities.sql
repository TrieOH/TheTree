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

-- name: ListEditionActivitiesAdmin :many
SELECT *
FROM activities
WHERE edition_id = $1;

-- name: GetAttendanceRecordByID :one
SELECT *
FROM attendance_records
WHERE id = $1;

-- name: ListActivityAttendanceRecords :many
SELECT *
FROM attendance_records
WHERE activity_id = $1;

-- name: GetUserActivityAttendanceRecords :many
SELECT *
FROM attendance_records
WHERE activity_id = $1 AND user_id = $2;

-- name: GetActiveUserActivityAttendanceRecords :one
SELECT *
FROM attendance_records
WHERE activity_id = $1 AND user_id = $2 AND status != 'cancelled';

-- name: IsUserRegistered :one
SELECT EXISTS (
    SELECT 1
    FROM attendance_records
    WHERE activity_id = $1
      AND user_id = $2
      AND status != 'cancelled'
) AS is_registered;

-- name: RegisterToActivity :one
INSERT INTO attendance_records (activity_id, user_id, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UnregisterFromActivity :exec
UPDATE attendance_records
SET
    status = 'cancelled',
    cancelled_at = now(),
    updated_at = now()
WHERE activity_id = $1 AND user_id = $2 AND status = 'registered';

-- name: MarkAttendanceRecordStatus :exec
UPDATE attendance_records
SET
    status = $1,
    updated_at = now(),
    scanned_by = $2
WHERE id = $3;

-- name: AttendanceExport :many
SELECT
    ar.id              AS attendance_id,
    ar.user_id,
    ar.status          AS attendance_status,
    ar.checked_in_at,
    ar.cancelled_at,
    ar.scanned_by,
    ar.created_at,
    a.id               AS activity_id,
    a.title            AS activity_title,
    a.location         AS activity_location,
    a.status           AS activity_status,
    a.starts_at        AS activity_starts_at,
    a.ends_at          AS activity_ends_at,
    a.difficulty
FROM attendance_records ar
JOIN activities a ON a.id = ar.activity_id
WHERE a.edition_id = @edition_id
  AND ar.deleted_at IS NULL
  AND a.deleted_at  IS NULL
  AND (sqlc.narg(activity_ids)::uuid[]                     IS NULL OR a.id           = ANY(sqlc.narg(activity_ids)::uuid[]))
  AND (sqlc.narg(attendance_statuses)::attendance_status[] IS NULL OR ar.status      = ANY(sqlc.narg(attendance_statuses)::attendance_status[]))
  AND (sqlc.narg(activity_statuses)::activity_status[]     IS NULL OR a.status       = ANY(sqlc.narg(activity_statuses)::activity_status[]))
  AND (sqlc.narg(difficulties)::difficulty_level[]         IS NULL OR a.difficulty   = ANY(sqlc.narg(difficulties)::difficulty_level[]))
  AND (sqlc.narg(date_from)::timestamptz                   IS NULL OR ar.created_at >= sqlc.narg(date_from))
  AND (sqlc.narg(date_to)::timestamptz                     IS NULL OR ar.created_at <= sqlc.narg(date_to))
ORDER BY a.starts_at ASC, ar.created_at ASC;

-- name: WaitlistAppend :one
--INSERT INTO waitlist_entries (activity_id, user_id, position)
--SELECT $1, $2, COALESCE(MAX(position) + 1, 0)
--FROM waitlist_entries
--WHERE activity_id = $1
--RETURNING *;

-- name: WaitlistGetFirst :one
--SELECT * FROM waitlist_entries
--WHERE activity_id = $1
--ORDER BY position
--LIMIT 1;

-- name: WaitlistRemove :exec
--DELETE FROM waitlist_entries WHERE id = $1;

-- name: WaitlistGetOrdered :many
--SELECT
--    w.*,
--    ROW_NUMBER() OVER (ORDER BY w.position) as queue_position
--FROM waitlist_entries w
--WHERE activity_id = $1
--ORDER BY w.position;