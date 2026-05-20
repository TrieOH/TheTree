package ports

import (
	"IdentityX/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, toCreate models.Session) (*models.Session, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*models.Session, error)
	GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*models.Session, error)
	List(ctx context.Context, userID uuid.UUID, userType models.UserType) ([]models.Session, error)
	Update(ctx context.Context, toUpdate models.Session, entityID uuid.UUID, userType models.UserType) error
	RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*models.Session, error)
	MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, userType models.UserType) (*models.Session, error)
	MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error
	MarkRevokedByFilter(ctx context.Context, filter models.Filter) (int, error)
}
