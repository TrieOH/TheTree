package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Capability) (*models.Capability, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
	defer span.End()
	capability, err := database.Queries(ctx, repo.q).CreateCapability(ctx, sqlc.CreateCapabilityParams{
		ProjectID: toCreate.ProjectID,
		Resource:  toCreate.Resource,
		Action:    toCreate.Action,
		CreatedBy: toCreate.CreatedBy,
	})
	return new(mapCapability(capability)), repo.dbe(err)
}
