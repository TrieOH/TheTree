package outbounds

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenReuseListRepository interface {
	Append(ctx context.Context, jit, userID uuid.UUID, expiresAt time.Time) error
	Exists(ctx context.Context, jit, userID uuid.UUID) (bool, error)
	ClearExpired(ctx context.Context) error
}
