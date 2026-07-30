-- name: CreateCertificationTemplate :one
INSERT INTO certification_templates (edition_id, kind, name, description, design_data)
VALUES (@edition_id, @kind, @name, @description, @design_data)
RETURNING *;

-- name: GetCertificationTemplateByID :one
SELECT *
FROM certification_templates
WHERE id = @id
  AND deleted_at IS NULL;

-- name: ListCertificationTemplatesByEdition :many
SELECT *
FROM certification_templates
WHERE edition_id = @edition_id
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: ListCertificationTemplatesByEditionForEmission :many
SELECT *
FROM certification_templates
WHERE edition_id = @edition_id
ORDER BY created_at;

-- name: UpdateCertificationTemplate :one
UPDATE certification_templates
SET
    kind        = @kind,
    name        = @name,
    description = @description,
    design_data = @design_data,
    updated_at  = now()
WHERE id = @id
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteCertificationTemplate :exec
UPDATE certification_templates
SET
    deleted_at = now(),
    updated_at = now()
WHERE id = @id
  AND deleted_at IS NULL;

-- name: CreateCertificationTemplateProgram :exec
INSERT INTO certification_template_programs (template_id, program_id)
VALUES (@template_id, @program_id);

-- name: DeleteCertificationTemplatePrograms :exec
DELETE FROM certification_template_programs
WHERE template_id = @template_id;

-- name: ListCertificationTemplateProgramsByTemplate :many
SELECT *
FROM certification_template_programs
WHERE template_id = @template_id;

-- name: CreateCertification :one
INSERT INTO certifications (edition_id, template_id, registration_id, user_id, program_id, verification_hash)
VALUES (@edition_id, @template_id, @registration_id, @user_id, @program_id, @verification_hash)
RETURNING *;

-- name: GetCertificationByID :one
SELECT *
FROM certifications
WHERE id = @id;

-- name: GetCertificationByHash :one
SELECT *
FROM certifications
WHERE verification_hash = @verification_hash;

-- name: ListCertificationsByUser :many
SELECT *
FROM certifications
WHERE user_id = @user_id
ORDER BY issued_at DESC;

-- name: ListCertificationsByEdition :many
SELECT *
FROM certifications
WHERE edition_id = @edition_id
ORDER BY user_id, issued_at DESC;

-- name: HasCertForProgram :one
SELECT EXISTS (
    SELECT 1
    FROM certifications
    WHERE user_id = @user_id
      AND program_id = @program_id
) AS exists;

-- name: HasCertForRegistration :one
SELECT EXISTS (
    SELECT 1
    FROM certifications
    WHERE registration_id = @registration_id
      AND template_id IS NOT DISTINCT FROM @template_id
) AS exists;

-- name: InvalidateCertification :exec
UPDATE certifications
SET
    valid          = false,
    invalid_reason = @invalid_reason,
    updated_at     = now()
WHERE id = @id;

-- name: MarkCertificationEmailSent :exec
UPDATE certifications
SET
    email_sent = true,
    updated_at = now()
WHERE id = @id;

-- name: RecordCertEmissionError :exec
INSERT INTO cert_emission_errors (edition_id, user_id, template_id, program_id, error_message)
VALUES (@edition_id, @user_id, @template_id, @program_id, @error_message);

-- name: ListCertEmissionErrorsByEdition :many
SELECT *
FROM cert_emission_errors
WHERE edition_id = @edition_id
ORDER BY created_at DESC;

-- name: ListDistinctRegistrationsByEdition :many
SELECT
    r.attendee_user_id AS user_id,
    r.id AS registration_id,
    r.attendee_email,
    r.attendee_name
FROM registrations r
WHERE r.edition_id = @edition_id
  AND r.status = 'confirmed'
  AND r.deleted_at IS NULL
GROUP BY r.attendee_user_id, r.id, r.attendee_email, r.attendee_name;

-- name: ListDistinctParticipantsByProgram :many
SELECT
    r.attendee_user_id AS user_id,
    r.id AS registration_id,
    r.attendee_email,
    r.attendee_name
FROM program_participations pp
JOIN registrations r ON r.id = pp.registration_id AND r.deleted_at IS NULL
JOIN program_occurrences po ON po.id = pp.occurrence_id AND po.deleted_at IS NULL
WHERE po.program_id = @program_id
GROUP BY r.attendee_user_id, r.id, r.attendee_email, r.attendee_name;
