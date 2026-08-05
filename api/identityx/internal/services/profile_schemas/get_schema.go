package profile_schemas

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

// GetSchema returns the active profile schema for the scope. Public read:
// schemas shape public profiles, so anyone may fetch them.
func (o *Operations) GetSchema(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetSchema")
	defer span.End()

	return o.schemas.Get(ctx, projectID)
}
