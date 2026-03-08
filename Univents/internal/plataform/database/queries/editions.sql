-- name: CreateEdition :one
INSERT INTO editions (
    id,
    event_id,
    goauth_scope_id,
    type,
    edition_name,
    -- name_template,
    tagline,
    description,
    status,
    -- visible_from,
    registration_opens_at,
    registration_closes_at,
    monetary_type,
    -- has_tickets,
    -- has_ticket_tiers,
    -- has_capacity,
    -- capacity,
    -- remaining_capacity,
    -- has_products,
    -- has_bundles,
    -- has_tokens,
    -- max_tokens_per_user,
    starts_at,
    ends_at,
    timezone,
    location_name,
    location_address,
    logo_url,
    banner_url,
    -- has_gallery,
    -- gallery_urls,
    -- has_schedule,
    -- schedule_id,
    -- has_activities,
    -- has_activity_interest_marking,
    -- has_paid_activities,
    -- has_interest_list,
    -- interest_list_opens_at,
    -- has_checkout,
    contact_email,
    contact_phone,
    organizer_name,
    created_by,
    created_at,
    updated_at
)
VALUES (
    $1,  
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18,
    $19,
    $20,
    $21,
    $22,
    now(),
    now()
)
RETURNING *;

-- name: ListEditions :many
SELECT *
FROM editions
WHERE event_id = $1 AND status != 'draft';

-- name: ListEditionsAdmin :many
SELECT *
FROM editions
WHERE event_id = $1;

-- name: GetEditionByID :one
SELECT *
FROM editions
WHERE id = $1;

-- name: AnnounceEdition :exec
UPDATE editions
SET status = 'announced'
WHERE id = $1 AND status = 'draft';

-- name: OpenEditionRegistrations :exec
UPDATE editions
SET status = 'open'
WHERE id = $1 AND status = 'announced';

-- name: StartEdition :exec
UPDATE editions
SET status = 'ongoing'
WHERE id = $1 AND status = 'open';

-- name: FinishEdition :exec
UPDATE editions
SET status = 'finished'
WHERE id = $1 AND status = 'ongoing';