-- name: CreateForm :one
INSERT INTO forms (namespace_id, created_by, owner_id, name, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFormByID :one
SELECT * FROM forms
WHERE id = $1;

-- name: GetFormMember :one
SELECT *
FROM form_members
WHERE user_id = $1 AND form_id = $2;

-- name: AddFormMember :exec
INSERT INTO form_members (user_id, form_id, role, added_at, added_by)
VALUES ($1, $2, $3, $4, $5);

-- name: RemoveFormMember :exec
DELETE FROM form_members
WHERE user_id = $1 AND form_id = $2;

-- name: ListDirectFormMembers :many
SELECT *
FROM form_members
WHERE form_id = $1;

-- name: ListMyForms :many
SELECT *
FROM forms
WHERE owner_id = $1;

-- name: ListJoinedForms :many
SELECT f.*
FROM form_members fm
INNER JOIN forms f
ON fm.form_id = f.id
WHERE fm.user_id = $1;

-- name: ListNamespaceForms :many
SELECT *
FROM forms
WHERE forms.namespace_id = $1;