package outbound

import (
	"GoAuth/internal/domain/session"
	"context"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, new session.Session) (*session.Session, error)
	GetById(ctx context.Context, sessionID uuid.UUID) (*session.Session, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (*session.Session, error)
	List(ctx context.Context, userID uuid.UUID) ([]session.Session, error)
	Update(ctx context.Context, updated session.Session) error
	DeleteByFilter(ctx context.Context, filter session.Filter) ([]session.Session, error)
}
