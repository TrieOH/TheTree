package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Seller, error) {
	ctx, span := repo.tracer.Start(ctx, "SellerRepo.GetByID")
	defer span.End()

	sqlcSeller, err := database.Queries(ctx, repo.q).GetSellerByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapSeller(sqlcSeller)), nil
}
