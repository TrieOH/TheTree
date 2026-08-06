package ports

import (
	"context"
	"encoding/json"
	"univents/models"

	"github.com/google/uuid"
)

type BadgeTemplateRepo interface {
	Create(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error)
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error)
	ListByEditionIDs(ctx context.Context, editionIDs []uuid.UUID) ([]models.BadgeTemplate, error)
	Update(ctx context.Context, id uuid.UUID, name string, designData json.RawMessage) (*models.BadgeTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// BadgeEmissionRepo persists badge emissions. The design is never stored —
// it is derived at read time from the edition's badge_templates.
type BadgeEmissionRepo interface {
	Upsert(ctx context.Context, emission *models.BadgeEmission) (*models.BadgeEmission, error)
	Revoke(ctx context.Context, editionID, userID uuid.UUID, origin models.BadgeEmissionOrigin, reason string) error
	RevokeByRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error
	MarkEmailSent(ctx context.Context, id uuid.UUID) error
	ListViewsByUser(ctx context.Context, userID uuid.UUID) ([]models.BadgeEmissionView, error)
	ListViewsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEmissionView, error)
}

// BadgeStaffOps is the trigger surface the events feature calls when a member
// is added or removed.
type BadgeStaffOps interface {
	AwardStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error
	RevokeStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error
}

// BadgeEditionOps is the trigger surface the editions feature calls when an
// edition is published.
type BadgeEditionOps interface {
	AwardStaffBadgesForEdition(ctx context.Context, editionID uuid.UUID) error
}
