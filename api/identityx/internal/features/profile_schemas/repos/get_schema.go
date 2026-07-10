package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (r *schemaRepo) Get(ctx context.Context, projectID *uuid.UUID) (*models.ProjectProfileSchema, error) {
	ctx, span := database.Span(ctx, r.tracer, "GetProfileSchema")
	defer span.End()
	result, err := database.Queries(ctx, r.q).GetProfileSchema(ctx, projectID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapProfileSchema(result)), nil
}
