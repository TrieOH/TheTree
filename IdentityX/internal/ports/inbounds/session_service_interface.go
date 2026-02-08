package inbounds

import (
	"context"

	"github.com/google/uuid"
)

type SessionService interface {
	List(ctx context.Context) ([]OutputSession, error)
	RevokeByID(ctx context.Context, sessionID uuid.UUID, accessToken string) error
	RevokeOthers(ctx context.Context, accessToken string) error
	RevokeAll(ctx context.Context) error
	Me(ctx context.Context, accessToken, refreshToken string) (*MeOutput, error)
}
