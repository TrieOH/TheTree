package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, toCreate contracts.Session) (*contracts.Session, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*contracts.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*contracts.Session, error)
	GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*contracts.Session, error)
	List(ctx context.Context, entityID uuid.UUID, identityType contracts.IdentityType) ([]contracts.Session, error)
	Update(ctx context.Context, toUpdate contracts.Session, entityID uuid.UUID, identityType contracts.IdentityType) error
	RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*contracts.Session, error)
	MarkRevokedByID(ctx context.Context, entityID uuid.UUID, sessionID uuid.UUID, identityType contracts.IdentityType) (*contracts.Session, error)
	MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error
	MarkRevokedByFilter(ctx context.Context, filter contracts.Filter) (int, error)
	CreateIdentity(ctx context.Context, identityType contracts.IdentityType, entityID uuid.UUID) (*contracts.Identity, error)
	GetIdentityByEntityIDAndType(ctx context.Context, entityID uuid.UUID, identityType contracts.IdentityType) (*contracts.Identity, error)
	GetIdentityByIDAndType(ctx context.Context, identityID uuid.UUID, identityType contracts.IdentityType) (*contracts.Identity, error)
}
