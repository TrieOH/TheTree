-- name: RegisterActor :one
INSERT INTO actors (project_id, auth_method, password_hash, email, type, metadata)
VALUES (
    @project_id,
    @auth_method,
    @password_hash,
    @email,
    @type,
    @metadata
) RETURNING *;

-- name: GetActorByEmail :one
SELECT *
FROM actors
WHERE email = @email
  AND project_id IS NOT DISTINCT FROM @project_id;

-- name: UpdateActorLastLoginAt :exec
UPDATE actors
SET last_login_at = NOW()
WHERE id = @actor_id;

-- name: HasAnyActor :one
SELECT EXISTS (SELECT 1 FROM actors LIMIT 1) AS exists;