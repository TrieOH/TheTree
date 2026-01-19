package outbounds

import (
	"GoAuth/internal/domain/session"
	"context"
	"time"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, toCreate session.Session) (*session.Session, error)
	GetByID(ctx context.Context, sessionID uuid.UUID) (*session.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*session.Session, error)
	GetByFamilyID(ctx context.Context, familyID uuid.UUID) (*session.Session, error)
	List(ctx context.Context, entityID uuid.UUID, identityType session.IdentityType) ([]session.Session, error)
	Update(ctx context.Context, toUpdate session.Session, entityID uuid.UUID, identityType session.IdentityType) error
	RotateToken(ctx context.Context, familyID uuid.UUID, newTokenID uuid.UUID, oldTokenID uuid.UUID, expiresAt time.Time) (*session.Session, error)
	MarkRevokedByID(ctx context.Context, entityID uuid.UUID, sessionID uuid.UUID, identityType session.IdentityType) (*session.Session, error)
	MarkRevokedByFamilyID(ctx context.Context, familyID uuid.UUID) error
	MarkRevokedByFilter(ctx context.Context, filter session.Filter) (int, error)
	CreateIdentity(ctx context.Context, identityType session.IdentityType, entityID uuid.UUID) (*session.Identity, error)
	GetIdentityByEntityIDAndType(ctx context.Context, entityID uuid.UUID, identityType session.IdentityType) (*session.Identity, error)
	GetIdentityByIDAndType(ctx context.Context, identityID uuid.UUID, identityType session.IdentityType) (*session.Identity, error)
}
