package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcOrg, err := database.Queries(ctx, repo.q).CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
		Slug:    toCreate.Slug,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOrganization(sqlcOrg)), nil
}
