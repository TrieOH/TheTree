-- name: RegisterActor :one
INSERT INTO actors (project_id, auth_method, password_hash, email, type, metadata, verified_at)
VALUES (
    @project_id,
    @auth_method,
    @password_hash,
    @email,
    @type,
    @metadata,
    @verified_at
) RETURNING *;

-- name: GetActorByEmail :one
SELECT *
FROM actors
WHERE email = @email
  AND project_id IS NOT DISTINCT FROM @project_id;

-- name: GetActorByID :one
SELECT *
FROM actors
WHERE id = @id;

-- name: UpdateActorLastLoginAt :exec
UPDATE actors
SET last_login_at = NOW()
WHERE id = @actor_id;

-- name: UpdateActorVerifiedAt :exec
UPDATE actors
SET verified_at = @verified_at,
    updated_at = NOW()
WHERE id = @actor_id;

-- name: UpdateActorPasswordHash :exec
UPDATE actors
SET password_hash = @password_hash,
    updated_at = NOW()
WHERE id = @actor_id;

-- name: HasAnyActor :one
SELECT EXISTS (SELECT 1 FROM actors LIMIT 1) AS exists;

-- name: GetExternalIdentityByProviderAndSubject :one
-- Scoped by the identity's actor project: a platform login (NULL project)
-- only sees platform identities, a project login only sees identities whose
-- actor lives in that project. Mirrors GetActorByEmail's scope check.
SELECT e.*
FROM actor_external_identities e
JOIN actors a ON a.id = e.actor_id
WHERE e.provider = @provider
  AND e.subject = @subject
  AND a.project_id IS NOT DISTINCT FROM @project_id
ORDER BY e.created_at
LIMIT 1;

-- name: CreateExternalIdentity :one
INSERT INTO actor_external_identities (actor_id, provider, subject, email, encrypted_access_token, encrypted_refresh_token, token_expires_at)
VALUES (@actor_id, @provider, @subject, @email, @encrypted_access_token, @encrypted_refresh_token, @token_expires_at)
    RETURNING *;

-- name: UpdateExternalIdentityTokens :one
-- Same scope rule as the read: only the row whose actor lives in the given
-- project scope is refreshed, so a platform login can never touch a
-- project's identity (and vice versa).
UPDATE actor_external_identities
SET encrypted_access_token = @encrypted_access_token,
    encrypted_refresh_token = @encrypted_refresh_token,
    token_expires_at = @token_expires_at,
    updated_at = NOW()
WHERE id = (
    SELECT e.id
    FROM actor_external_identities e
    JOIN actors a ON a.id = e.actor_id
    WHERE e.provider = @provider
      AND e.subject = @subject
      AND a.project_id IS NOT DISTINCT FROM @project_id
    ORDER BY e.created_at
    LIMIT 1
)
RETURNING *;

-- name: ListActorsFromProject :many
SELECT *
FROM actors
WHERE project_id = @project_id;