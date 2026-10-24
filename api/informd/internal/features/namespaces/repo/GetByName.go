package repo

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error) {
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByName(ctx, sqlc.GetNamespaceByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
