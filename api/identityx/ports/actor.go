package ports

import (
	"IdentityX/models"
	"context"
	"time"
	"uuid"
)

type ActorRepo interface {
	Register(ctx context.Context, toRegister models.Actor) (*models.Actor, error)
	// LockSetup holds the setup lock until the transaction ends; call it
	// inside one, then HasAny.
	LockSetup(ctx context.Context) error
	// HasAny reports whether any actor exists, which is what "set up" means.
	HasAny(ctx context.Context) (bool, error)
	GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Actor, error)
	List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error)
	GetProjectServiceAccount(ctx context.Context, id uuid.UUID) (*models.Actor, error)
	UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error
	SetVerifiedAt(ctx context.Context, actorID uuid.UUID, at time.Time) error
	UpdatePasswordHash(ctx context.Context, actorID uuid.UUID, hash string) error
}
