package outbounds

import (
	"GoAuth/internal/domain/user"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, email, password string) (*user.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	Verify(ctx context.Context, userID uuid.UUID) (bool, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error
}
