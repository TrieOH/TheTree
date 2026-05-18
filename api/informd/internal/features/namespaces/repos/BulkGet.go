package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "BulkGet")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).BulkGetNamespaces(ctx, ids)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForm, mapNamespace), nil
}
