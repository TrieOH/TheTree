package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	sqlcNamespaces, err := database.Queries(ctx, repo.q).ListOwnedNamespaces(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcNamespaces, mapNamespace), nil
}
