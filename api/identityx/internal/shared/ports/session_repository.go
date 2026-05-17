package ports

import (
	"IdentityX/contracts"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, toCreate contracts.Session) (*contracts.Session, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*contracts.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*contracts.Session, error)
	GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*contracts.Session, error)
	List(ctx context.Context, userID uuid.UUID, userType contracts.UserType) ([]contracts.Session, error)
	Update(ctx context.Context, toUpdate contracts.Session, entityID uuid.UUID, userType contracts.UserType) error
	RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*contracts.Session, error)
	MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, userType contracts.UserType) (*contracts.Session, error)
	MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error
	MarkRevokedByFilter(ctx context.Context, filter contracts.Filter) (int, error)
}
