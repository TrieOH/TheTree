package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	results, err := database.Queries(ctx, repo.q).ListProductsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapProduct), nil
}
