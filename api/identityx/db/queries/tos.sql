-- name: GetCurrentTos :one
SELECT *
FROM terms_of_service
WHERE project_id = @project_id;

-- name: InsertTos :one
INSERT INTO terms_of_service (project_id, version, content, effective_at)
VALUES (@project_id, @version, @content, @effective_at)
RETURNING *;

-- name: UpdateTosBumpVersion :one
UPDATE terms_of_service
SET version = version + 1,
    content = @content,
    effective_at = @effective_at,
    updated_at = NOW()
WHERE project_id = @project_id
RETURNING *;

-- name: InsertTosAcceptance :one
INSERT INTO tos_acceptances (actor_id, project_id, tos_version, content_hash, source)
VALUES (@actor_id, @project_id, @tos_version, @content_hash, @source)
ON CONFLICT (actor_id, project_id, tos_version)
DO UPDATE SET accepted_at = NOW()
RETURNING *;

-- The newest version the actor has any acceptance row for, 0 when none.
-- The continued-use stamp compares it against the current version.
-- name: GetLatestAcceptedTosVersion :one
SELECT COALESCE(MAX(tos_version), 0)::int AS latest_version
FROM tos_acceptances
WHERE actor_id = @actor_id
  AND project_id = @project_id;

-- name: ListTosAcceptances :many
SELECT ta.*, a.email AS actor_email
FROM tos_acceptances ta
JOIN actors a ON a.id = ta.actor_id
WHERE ta.project_id = @project_id
ORDER BY ta.accepted_at DESC;
