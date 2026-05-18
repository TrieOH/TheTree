package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByName")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByName(ctx, sqlc.GetNamespaceByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
