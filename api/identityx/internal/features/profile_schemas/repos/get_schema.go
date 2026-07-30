package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"

	"github.com/google/uuid"
)

func (r *Repo) Get(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := telemetry.StartSpan(ctx, "Get")
	defer span.End()

	result, err := database.Queries(ctx, r.q).GetProfileSchema(ctx, projectID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapProfileSchema(result)), nil
}
