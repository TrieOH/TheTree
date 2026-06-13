package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListMine(ctx context.Context, userID uuid.UUID) ([]models.Form, error) {
	ctx, span := repo.tracer.Start(ctx, "ListMine")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).ListMyForms(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
