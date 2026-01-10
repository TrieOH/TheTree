package outbound

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
	List(ctx context.Context, userID uuid.UUID) ([]session.Session, error)
	Update(ctx context.Context, toUpdate session.Session) error
	RotateToken(ctx context.Context, oldTokenID uuid.UUID, newTokenID uuid.UUID, expiresAt time.Time) (*session.Session, error)
	MarkRevokedByID(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*session.Session, error)
	MarkRevokedByFilter(ctx context.Context, filter session.Filter) (int, error)
}
