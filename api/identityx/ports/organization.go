package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type OrganizationRepo interface {
	Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	AddMember(ctx context.Context, toCreate models.OrganizationMember) error
	RemoveMember(ctx context.Context, actorID, orgID uuid.UUID) error
	GetMember(ctx context.Context, actorID, orgID uuid.UUID) (*models.OrganizationMember, error)
	ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error)
}
