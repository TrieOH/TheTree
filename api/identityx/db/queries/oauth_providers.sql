-- name: CreateProjectOAuthProvider :one
INSERT INTO project_oauth_providers (project_id, provider, client_id, encrypted_client_secret)
VALUES (@project_id, @provider, @client_id, @encrypted_client_secret)
RETURNING *;

-- name: ListProjectOAuthProviders :many
SELECT *
FROM project_oauth_providers
WHERE project_id = @project_id
ORDER BY created_at;

-- name: GetProjectOAuthProvider :one
SELECT *
FROM project_oauth_providers
WHERE id = @id;

-- name: GetProjectOAuthProviderByProjectAndProvider :one
SELECT *
FROM project_oauth_providers
WHERE project_id = @project_id
  AND provider = @provider;

-- name: UpdateProjectOAuthProviderClientID :one
UPDATE project_oauth_providers
SET client_id = @client_id,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateProjectOAuthProviderClientSecret :one
UPDATE project_oauth_providers
SET encrypted_client_secret = @encrypted_client_secret,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SetProjectOAuthProviderEnabled :one
UPDATE project_oauth_providers
SET enabled = @enabled,
    updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteProjectOAuthProvider :exec
DELETE FROM project_oauth_providers
WHERE id = @id;

-- name: CreateOAuthLoginState :one
INSERT INTO oauth_login_states (state, provider, project_id, expires_at)
VALUES (@state, @provider, @project_id, @expires_at)
RETURNING *;

-- name: GetOAuthLoginState :one
SELECT *
FROM oauth_login_states
WHERE state = @state;

-- name: DeleteOAuthLoginState :exec
DELETE FROM oauth_login_states
WHERE id = @id;
