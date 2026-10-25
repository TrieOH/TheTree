package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error) {
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
