package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByForm(ctx context.Context, formID uuid.UUID) ([]models.Response, error) {
	rows, err := database.Queries(ctx, repo.q).ListResponsesByForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapResponse), nil
}
