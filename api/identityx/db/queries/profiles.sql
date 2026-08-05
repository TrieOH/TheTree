-- name: GetActorProfile :one
SELECT *
FROM actor_profiles
WHERE actor_id = @actor_id;

-- name: GetActorProfileByHandle :one
SELECT *
FROM actor_profiles
WHERE handle = @handle;

-- name: UpsertActorProfile :one
INSERT INTO actor_profiles (actor_id, profile, schema_version, outdated, handle, updated_at)
VALUES (@actor_id, @profile, @schema_version, @outdated, @handle, NOW())
ON CONFLICT (actor_id) DO UPDATE SET
    profile        = EXCLUDED.profile,
    schema_version = EXCLUDED.schema_version,
    outdated       = EXCLUDED.outdated,
    handle         = EXCLUDED.handle,
    updated_at     = NOW()
RETURNING *;

-- name: SetActorProfileMigrationState :one
UPDATE actor_profiles
SET schema_version = @schema_version,
    outdated       = @outdated,
    updated_at     = NOW()
WHERE actor_id = @actor_id
RETURNING *;

-- name: ListOutdatedProfiles :many
-- Outdated profiles for the scope; @project_id NULL means the platform
-- scope (actors with no project).
SELECT ap.*
FROM actor_profiles ap
JOIN actors a ON a.id = ap.actor_id
WHERE ap.outdated = true AND a.project_id IS NOT DISTINCT FROM @project_id
ORDER BY ap.updated_at DESC;

-- name: GetProfileSchema :one
SELECT *
FROM project_profile_schemas
WHERE project_id IS NOT DISTINCT FROM @project_id;

-- name: UpsertProfileSchema :one
-- Idempotent, versioned upsert: setting the same schema again keeps the
-- current version; a changed schema bumps the version and appends a row to
-- profile_schema_versions (history). Runs as a single statement, so the
-- version bump and its history row commit atomically.
WITH old_schema AS (
    SELECT version
    FROM project_profile_schemas
    WHERE project_profile_schemas.project_id IS NOT DISTINCT FROM @project_id
), upserted AS (
    INSERT INTO project_profile_schemas (project_id, schema, version, active, updated_at)
    VALUES (@project_id, @schema, 1, @active, NOW())
    ON CONFLICT (project_id) DO UPDATE SET
        schema     = EXCLUDED.schema,
        version    = CASE
            WHEN project_profile_schemas.schema = EXCLUDED.schema
                THEN project_profile_schemas.version
            ELSE project_profile_schemas.version + 1
        END,
        active     = EXCLUDED.active,
        updated_at = NOW()
    RETURNING *
), history AS (
    INSERT INTO profile_schema_versions (project_id, version, schema)
    SELECT upserted.project_id, upserted.version, upserted.schema
    FROM upserted
    LEFT JOIN old_schema ON TRUE
    WHERE upserted.version IS DISTINCT FROM COALESCE(old_schema.version, -1)
)
SELECT * FROM upserted;
