package ports

import (
	"IdentityX/models"
	"context"
	"uuid"
)

// ProfileSchemaRepo handles project-level profile schema management.
// When projectID is nil, the schema is the platform-wide default.
type ProfileSchemaRepo interface {
	Get(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error)
	Upsert(ctx context.Context, schema models.ProjectProfileSchema) (*models.ProjectProfileSchema, error)
}
