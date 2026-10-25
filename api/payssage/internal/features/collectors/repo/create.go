package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Collector) (*models.Collector, error) {
	sqlcCollector, err := database.Queries(ctx, repo.q).CreateCollector(ctx, sqlc.CreateCollectorParams{
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Provider:       toCreate.Provider,
		ProviderUserID: toCreate.ProviderUserID,
		Credentials:    toCreate.Credentials,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapCollector(sqlcCollector)), nil
}
