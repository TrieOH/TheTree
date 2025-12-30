package outbound

import (
	"GoAuth/internal/domain/revoked_refreshes"
	"context"
	"time"

	"github.com/google/uuid"
)

type RevokedRefreshTokenRepository interface {
	Revoke(ctx context.Context, blacklist revoked_refreshes.RevokedRefreshToken) error
	RevokeMany(ctx context.Context, tokenIDs []uuid.UUID, expiresAts []time.Time) error
	GetByID(ctx context.Context, BlacklistID uuid.UUID) (*revoked_refreshes.RevokedRefreshToken, error)
	Delete(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	IsRevoked(ctx context.Context, tokenID uuid.UUID) (bool, error)
}
