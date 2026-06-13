package repos

import (
	"context"

	"Informd/internal/database/sqlc"
	"Informd/models"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).CreateNamespace(ctx, sqlc.CreateNamespaceParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
