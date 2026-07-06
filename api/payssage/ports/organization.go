package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type OrganizationRepo interface {
	Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error)
	GetByID(ctx context.Context, orgID uuid.UUID) (*models.Organization, error)
	ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error)
	AddMember(ctx context.Context, toAdd models.OrganizationMember) error
	RemoveMember(ctx context.Context, memberID, orgID uuid.UUID) error
	GetMember(ctx context.Context, memberID, orgID uuid.UUID) (*models.OrganizationMember, error)
	ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error)
}
