package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "ListJoined")
	defer span.End()
	sqlcNamespaces, err := database.Queries(ctx, repo.q).ListJoinedNamespaces(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcNamespaces, mapNamespace), nil
}
