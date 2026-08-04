package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// ProfileRepo handles actor profile CRUD. The profile is a JSONB document
// that downstream project owners shape via a JSON Schema. Each profile
// carries the schema version it was validated against plus an outdated flag.
type ProfileRepo interface {
	Get(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error)
	// GetByHandle returns the profile whose handle matches, or an error when
	// no profile has that handle. Handles are unique when present.
	GetByHandle(ctx context.Context, handle string) (*models.ActorProfile, error)
	Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error)
	// SetMigrationState bumps the schema version pointer or toggles the
	// outdated flag without touching the profile document itself.
	SetMigrationState(ctx context.Context, actorID uuid.UUID, schemaVersion int, outdated bool) (*models.ActorProfile, error)
	// ListOutdated returns profiles flagged as outdated. A nil projectID
	// means the platform scope (actors with no project).
	ListOutdated(ctx context.Context, projectID *uuid.UUID) ([]models.ActorProfile, error)
}
