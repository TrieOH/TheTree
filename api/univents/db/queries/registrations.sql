-- name: GetRegistrationByID :one
SELECT *
FROM registrations
WHERE id = @id;
