package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ActorRepo interface {
	Register(ctx context.Context, toRegister models.Actor) (*models.Actor, error)
	GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Actor, error)
	UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error
}
