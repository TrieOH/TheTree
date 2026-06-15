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

-- name: GetActorByID :one
SELECT *
FROM actors
WHERE id = @id;

-- name: UpdateActorLastLoginAt :exec
UPDATE actors
SET last_login_at = NOW()
WHERE id = @actor_id;

-- name: HasAnyActor :one
SELECT EXISTS (SELECT 1 FROM actors LIMIT 1) AS exists;

-- name: GetExternalIdentityByProviderAndSubject :one
SELECT * FROM actor_external_identities
WHERE provider = @provider
  AND subject = @subject;

-- name: CreateExternalIdentity :one
INSERT INTO actor_external_identities (actor_id, provider, subject, email, encrypted_access_token, encrypted_refresh_token, token_expires_at)
VALUES (@actor_id, @provider, @subject, @email, @encrypted_access_token, @encrypted_refresh_token, @token_expires_at)
    RETURNING *;

-- name: UpdateExternalIdentityTokens :one
UPDATE actor_external_identities
SET encrypted_access_token = @encrypted_access_token,
    encrypted_refresh_token = @encrypted_refresh_token,
    token_expires_at = @token_expires_at,
    updated_at = NOW()
WHERE provider = @provider
  AND subject = @subject
    RETURNING *;