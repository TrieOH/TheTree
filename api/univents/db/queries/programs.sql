-- name: CreateProgram :one
INSERT INTO programs (edition_id, kind, name, description, min_access_level, staff_only, price)
VALUES (@edition_id, @kind, @name, @description, @min_access_level, @staff_only, @price)
RETURNING *;

-- name: GetProgramByID :one
SELECT *
FROM programs
WHERE id = @id
  AND deleted_at IS NULL;

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
