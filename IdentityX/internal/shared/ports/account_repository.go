package ports

import (
	"context"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Verify(ctx context.Context, userID uuid.UUID) (bool, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error
}
