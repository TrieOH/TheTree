package repo

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error) {
	sqlcProject, err := database.Queries(ctx, repo.q).CreateNamespace(ctx, sqlc.CreateNamespaceParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
