-- name: GetActorProfile :one
SELECT *
FROM actor_profiles
WHERE actor_id = @actor_id;

-- name: UpsertActorProfile :one
INSERT INTO actor_profiles (actor_id, profile, updated_at)
VALUES (@actor_id, @profile, NOW())
ON CONFLICT (actor_id) DO UPDATE SET
    profile   = EXCLUDED.profile,
    updated_at = NOW()
RETURNING *;

-- name: GetProfileSchema :one
SELECT *
FROM project_profile_schemas
WHERE project_id IS NOT DISTINCT FROM @project_id;

-- name: UpsertProfileSchema :one
INSERT INTO project_profile_schemas (project_id, schema, version, active, updated_at)
VALUES (@project_id, @schema, 1, @active, NOW())
ON CONFLICT (project_id) DO UPDATE SET
    schema     = EXCLUDED.schema,
    version    = project_profile_schemas.version + 1,
    active     = EXCLUDED.active,
    updated_at = NOW()
RETURNING *;
