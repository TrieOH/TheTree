package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListMineArchived(ctx context.Context, userID uuid.UUID) ([]models.Form, error) {
	ctx, span := repo.tracer.Start(ctx, "ListMineArchived")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).ListMyArchivedForms(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
