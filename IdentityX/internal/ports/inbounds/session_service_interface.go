package inbounds

import (
	"GoAuth/internal/application/authz"
	"context"
)

type SessionService interface {
	List(ctx context.Context) ([]OutputSession, error)
	RevokeByID(ctx context.Context, sessionID string) error
	RevokeOthers(ctx context.Context) error
	RevokeAll(ctx context.Context) error
	Me(ctx context.Context) (*authz.Principal, error)
}
