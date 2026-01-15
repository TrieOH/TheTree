package inbounds

import (
	"context"

	"github.com/google/uuid"
)

type SessionService interface {
	List(ctx context.Context) ([]OutputSession, error)
	RevokeByID(ctx context.Context, sessionID uuid.UUID) error
	RevokeOthers(ctx context.Context) error
	RevokeAll(ctx context.Context) error
	Me(ctx context.Context) (*PrincipalOutput, error)
}
