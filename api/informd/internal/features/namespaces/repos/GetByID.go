package repos

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByID")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
