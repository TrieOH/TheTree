-- name: GivePlatformRole :one
INSERT INTO platform_roles (actor_id, role, metadata)
VALUES (
    @actor_id,
    @role,
    @metadata
) RETURNING *;

-- name: GetPlatformRoleByActor :one
SELECT role FROM platform_roles
WHERE actor_id = @actor_id;
