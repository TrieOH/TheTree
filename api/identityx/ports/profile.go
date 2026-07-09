package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

// ProfileRepo handles actor profile CRUD. The profile is a JSONB document
// that downstream project owners shape via a JSON Schema.
type ProfileRepo interface {
	Get(ctx context.Context, actorID uuid.UUID) (*models.ActorProfile, error)
	Upsert(ctx context.Context, profile models.ActorProfile) (*models.ActorProfile, error)
}

// ProfileSchemaRepo handles project-level profile schema management.
// When projectID is nil, the schema is the platform-wide default.
type ProfileSchemaRepo interface {
	Get(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error)
	Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error)
}
