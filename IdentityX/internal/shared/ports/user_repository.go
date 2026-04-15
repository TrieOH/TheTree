package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, email, password string) (*contracts.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*contracts.User, error)
	GetUserByEmail(ctx context.Context, email string) (*contracts.User, error)
	Verify(ctx context.Context, userID uuid.UUID) (bool, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error
}
