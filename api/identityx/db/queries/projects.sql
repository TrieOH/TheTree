-- name: ListProjects :many
SELECT * FROM projects;

-- name: CreateProject :one
INSERT INTO projects (organization_id, owner_id, brand_slug, name, domain, metadata)
VALUES (@organization_id, @owner_id, @brand_slug, @name,@domain, @metadata)
RETURNING *;

-- name: GetProjectByID :one
SELECT *
FROM projects
WHERE id = @id;

-- name: ListProjectsByOrganizationID :many
SELECT *
FROM projects
WHERE organization_id = @organization_id
  AND deleted_at IS NULL
  AND organization_id IS NOT NULL;

-- name: ListOwnedProjects :many
SELECT *
FROM projects
WHERE owner_id = @owner_id
  AND organization_id IS NULL;

-- name: ListJoinedProjects :many
SELECT p.*
FROM project_members pm
INNER JOIN projects p
ON pm.project_id = p.id
WHERE pm.actor_id = @actor_id
  AND p.owner_id != @actor_id;

-- name: AddProjectMember :exec
INSERT INTO project_members (project_id, actor_id, role, metadata, joined_at)
VALUES (@project_id, @actor_id, @role, @metadata, NOW())
RETURNING *;

-- name: RemoveProjectMember :exec
DELETE FROM project_members
WHERE project_id = @project_id
  AND actor_id = @actor_id;

-- name: GetProjectMemberByID :one
SELECT * FROM project_members
WHERE project_id = @project_id
  AND actor_id = @actor_id;

-- name: ListProjectMembers :many
SELECT * FROM project_members
WHERE project_id = @project_id;

-- name: GetProjectServiceAccount :one
SELECT * FROM actors
WHERE project_id = @project_id
  AND type = 'service';