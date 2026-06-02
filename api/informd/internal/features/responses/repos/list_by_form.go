package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) ListByForm(ctx context.Context, formID uuid.UUID) ([]models.Response, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponseRepo.ListByForm")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListResponsesByForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapResponse), nil
}
