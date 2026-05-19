package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "ListOwned")
	defer span.End()
	sqlcNamespaces, err := database.Queries(ctx, repo.q).ListOwnedNamespaces(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcNamespaces, mapNamespace), nil
}
