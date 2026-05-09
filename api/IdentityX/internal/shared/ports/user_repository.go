package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, email, password string, projectID *uuid.UUID, userType contracts.UserType) (*contracts.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*contracts.User, error)
	GetUserByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*contracts.User, error)
	ListFromProject(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error)
	GetByIDFromProject(ctx context.Context, userID, projectID uuid.UUID) (*contracts.User, error)
}
