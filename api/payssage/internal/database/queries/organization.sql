-- name: CreateOrganization :one
INSERT INTO organizations (owner_id, name, slug)
VALUES (@owner_id, @name, @slug)
RETURNING *;

-- name: GetOrganizationByID :one
SELECT *
FROM organizations
WHERE id = @id;

-- name: ListOwnedOrganizations :many
SELECT *
FROM organizations
WHERE owner_id = @owner_id;

-- name: ListJoinedOrganizations :many
SELECT o.*
FROM org_members om
INNER JOIN organizations o
   ON om.organization_id = o.id
WHERE om.member_id = @member_id
  AND o.owner_id != @member_id;

-- name: GetOrganizationMember :one
SELECT *
FROM org_members
WHERE member_id = @member_id
  AND organization_id = @organization_id;

-- name: AddOrganizationMember :exec
INSERT INTO org_members (member_id, organization_id, role, joined_at)
VALUES (@member_id, @organization_id, @role, NOW());

-- name: RemoveOrganizationMember :exec
DELETE FROM org_members
WHERE member_id = @member_id
  AND organization_id = @organization_id;

-- name: ListOrganizationMembers :many
SELECT *
FROM org_members
WHERE organization_id = @organization_id;