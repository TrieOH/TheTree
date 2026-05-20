package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(ctx context.Context, email, password string, projectID *uuid.UUID, userType models.UserType) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.User, error)
	ListFromProject(ctx context.Context, projectID uuid.UUID) ([]models.User, error)
	GetByIDFromProject(ctx context.Context, userID, projectID uuid.UUID) (*models.User, error)
}
