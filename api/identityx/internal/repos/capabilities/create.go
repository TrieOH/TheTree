package capabilities

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Capability) (*models.Capability, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	capability, err := database.Queries(ctx, repo.q).CreateCapability(ctx, sqlc.CreateCapabilityParams{
		ProjectID: toCreate.ProjectID,
		Resource:  toCreate.Resource,
		Action:    toCreate.Action,
		CreatedBy: toCreate.CreatedBy,
	})
	return new(mapCapability(capability)), repo.dbe(err)
}
