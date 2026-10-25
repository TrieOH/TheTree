package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Capability) (*models.Capability, error) {
	capability, err := database.Queries(ctx, repo.q).CreateCapability(ctx, sqlc.CreateCapabilityParams{
		ProjectID: toCreate.ProjectID,
		Resource:  toCreate.Resource,
		Action:    toCreate.Action,
		CreatedBy: toCreate.CreatedBy,
	})
	return new(mapCapability(capability)), repo.dbe(err)
}
