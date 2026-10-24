package repo

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error) {
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
