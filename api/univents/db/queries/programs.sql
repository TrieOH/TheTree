-- name: CreateProgramParticipation :one
INSERT INTO program_participations (edition_id, occurrence_id, registration_id, status)
VALUES (@edition_id, @occurrence_id, @registration_id, @status)
RETURNING *;

-- name: UpsertParticipationAttended :one
-- The checkpoint check-in upsert: staff marks an attendee present by
-- creating the attended participation on first scan (the row is anchored
-- to the attendee's edition-ticket registration — checkpoints have no
-- sign-up of their own) or flipping an existing live row back to
-- attended (idempotent re-scan). The partial unique index
-- (uniq_program_participations_active_per_occurrence_attendee) is the
-- conflict target: a cancelled row is not in the index, so re-checking-in
-- after a cancel inserts a fresh row (append-only ledger, like register).
INSERT INTO program_participations (edition_id, occurrence_id, registration_id, status)
VALUES (@edition_id, @occurrence_id, @registration_id, 'attended')
ON CONFLICT (occurrence_id, registration_id) WHERE status IN ('registered', 'attended', 'no_show')
DO UPDATE SET
    status     = 'attended',
    updated_at = now()
RETURNING *;

-- name: GetProgramParticipationByID :one
SELECT *
FROM program_participations
WHERE id = @id;

-- name: UpdateProgramParticipationStatus :one
UPDATE program_participations
SET
    status     = @status,
    updated_at = now()
WHERE id = @id
RETURNING *;

-- name: UpdateProgramParticipationStatusIfRegistered :one
-- The de-registration guard: flips only rows still in 'registered'. 0 rows
-- returned = the caller raced a state change or the spot is locked
-- (attended/no_show) — the caller maps it to 409.
UPDATE program_participations
SET
    status     = @status,
    updated_at = now()
WHERE id = @id
  AND status = 'registered'
RETURNING *;

-- name: UpdateProgramParticipationStatusIfNotCancelled :one
-- The mark-attended guard: flips any non-cancelled row — idempotent
-- attended→attended re-mark, allows the no_show→attended staff correction.
-- 0 rows returned = cancelled (refunded purchase / dropped out) — the
-- caller maps it to 409.
UPDATE program_participations
SET
    status     = @status,
    updated_at = now()
WHERE id = @id
  AND status <> 'cancelled'
RETURNING *;

-- name: GetActiveProgramParticipationByOccurrenceAndRegistration :one
-- The caller's live participation in an occurrence (register pre-check →
-- 409 already registered; de-register lookup → 404 not registered).
SELECT *
FROM program_participations
WHERE occurrence_id = @occurrence_id
  AND registration_id = @registration_id
  AND status IN ('registered', 'attended', 'no_show')
LIMIT 1;

-- name: CountActiveProgramParticipationsByOccurrence :one
-- Occupancy of an occurrence: every non-cancelled participation holds a
-- slot (paid spots materialize a participation too, so this covers both
-- the checkout and self-service paths). Runs inside the register tx under
-- the occurrence row lock.
SELECT COUNT(*)::BIGINT
FROM program_participations
WHERE occurrence_id = @occurrence_id
  AND status IN ('registered', 'attended', 'no_show');

-- name: ListActiveProgramParticipationsByEditionAndRegistration :many
-- The "my activities" read: the caller's live participations in an edition
-- joined with their program and occurrence, so the front can render the
-- schedule without extra round-trips. Cancelled rows are history — never
-- shown here.
SELECT
    pp.id                          AS participation_id,
    pp.edition_id                  AS edition_id,
    pp.occurrence_id               AS occurrence_id,
    pp.registration_id             AS registration_id,
    pp.status                      AS status,
    pp.created_at                  AS participation_created_at,
    pp.updated_at                  AS participation_updated_at,
    p.id                           AS program_id,
    p.kind                         AS program_kind,
    p.name                         AS program_name,
    p.description                  AS program_description,
    p.min_access_level             AS program_min_access_level,
    p.staff_only                   AS program_staff_only,
    p.price                        AS program_price,
    p.banner_url                   AS program_banner_url,
    po.starts_at                   AS occurrence_starts_at,
    po.ends_at                     AS occurrence_ends_at,
    po.max_capacity                AS occurrence_max_capacity
FROM program_participations pp
JOIN program_occurrences po ON po.id = pp.occurrence_id AND po.deleted_at IS NULL
JOIN programs p ON p.id = po.program_id AND p.deleted_at IS NULL
WHERE pp.edition_id = @edition_id
  AND pp.registration_id = @registration_id
  AND pp.status IN ('registered', 'attended', 'no_show')
ORDER BY po.starts_at;

-- name: ListProgramParticipationsByOccurrence :many
-- The staff marking surface: who is expected at an occurrence, with the
-- attendee identity from their registration. Cancelled rows are history —
-- staff mark the people who are coming, not the ones who dropped out.
SELECT
    pp.id                          AS participation_id,
    pp.occurrence_id               AS occurrence_id,
    pp.registration_id             AS registration_id,
    pp.status                      AS status,
    pp.created_at                  AS participation_created_at,
    r.attendee_user_id             AS attendee_user_id,
    r.attendee_email               AS attendee_email,
    r.attendee_name                AS attendee_name
FROM program_participations pp
JOIN registrations r ON r.id = pp.registration_id AND r.deleted_at IS NULL
WHERE pp.occurrence_id = @occurrence_id
  AND pp.status IN ('registered', 'attended', 'no_show')
ORDER BY pp.created_at;

-- name: CreateProgram :one
INSERT INTO programs (edition_id, kind, name, description, min_access_level, staff_only, price)
VALUES (@edition_id, @kind, @name, @description, @min_access_level, @staff_only, @price)
RETURNING *;

-- name: GetProgramByID :one
SELECT *
FROM programs
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetProgramByIDForUpdate :one
-- Row-lock variant for the checkout tx (split 7): serializes concurrent
-- checkouts on the same program (price read) before availability is
-- checked.
SELECT *
FROM programs
WHERE id = @id
  AND deleted_at IS NULL
FOR UPDATE;

-- name: ListProgramsByEdition :many
SELECT *
FROM programs
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: PatchProgram :one
UPDATE programs
SET
    kind             = @kind,
    name             = @name,
    description      = @description,
    min_access_level = @min_access_level,
    staff_only       = @staff_only,
    price            = @price,
    banner_url       = @banner_url,
    updated_at       = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProgram :one
UPDATE programs
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: CreateProgramOccurrence :one
INSERT INTO program_occurrences (program_id, edition_id, starts_at, ends_at, max_capacity)
VALUES (@program_id, @edition_id, @starts_at, @ends_at, @max_capacity)
RETURNING *;

-- name: GetProgramOccurrenceByID :one
SELECT *
FROM program_occurrences
WHERE id = @id
  AND deleted_at IS NULL;

-- name: GetProgramOccurrenceByIDForUpdate :one
-- Row-lock variant for the checkout tx (split 7): serializes concurrent
-- checkouts on the same occurrence (capacity) before availability is
-- checked.
SELECT *
FROM program_occurrences
WHERE id = @id
  AND deleted_at IS NULL
FOR UPDATE;

-- name: ListProgramOccurrencesByProgram :many
SELECT *
FROM program_occurrences
WHERE program_id = @program_id
  AND deleted_at IS NULL
ORDER BY starts_at;

-- name: ListProgramOccurrencesByEdition :many
SELECT *
FROM program_occurrences
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY starts_at;

-- name: PatchProgramOccurrence :one
UPDATE program_occurrences
SET
    starts_at    = @starts_at,
    ends_at      = @ends_at,
    max_capacity = @max_capacity,
    updated_at   = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProgramOccurrence :one
UPDATE program_occurrences
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;
