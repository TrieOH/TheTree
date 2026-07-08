-- name: CreateCertificationTemplate :one
INSERT INTO certification_templates (edition_id, title, data, url)
VALUES (@edition_id, @title, @data, @url)
RETURNING *;

-- name: ListCertificationTemplates :many
SELECT *
FROM certification_templates
WHERE edition_id = @edition_id;

-- name: GetCertificationTemplateByID :one
SELECT *
FROM certification_templates
WHERE id = @id
  AND edition_id = @edition_id;

-- name: Certify :one
INSERT INTO certifications (user_id, target_id, target_type)
VALUES (@user_id, @target_id, @target_type)
RETURNING *;

-- name: ListUserCertifications :many
SELECT *
FROM certifications
WHERE user_id = @user_id;

-- name: ListTargetCertifications :many
SELECT *
FROM certifications
WHERE target_type = @target_type
  AND target_id = @target_id;

-- name: GetCertificationByID :one
SELECT *
FROM certifications
WHERE id = @id;

-- name: SetActivityCertificationTemplate :exec
UPDATE activities
SET
    certification_template_id = @certification_template_id
WHERE id = @id;

-- name: SetEditionCertificationTemplate :exec
UPDATE editions
SET
    certification_template_id = @certification_template_id
WHERE id = @id;