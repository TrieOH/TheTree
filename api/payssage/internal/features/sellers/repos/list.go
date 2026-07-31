package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"
)

func (repo *Repo) List(ctx context.Context) ([]models.Seller, error) {
	ctx, span := repo.tracer.Start(ctx, "SellerRepo.List")
	defer span.End()

	sqlcSellers, err := database.Queries(ctx, repo.q).ListSellers(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcSellers, mapSeller), nil
}
