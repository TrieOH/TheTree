package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"
)

func (repo *Repo) List(ctx context.Context) ([]models.Seller, error) {
	sqlcSellers, err := database.Queries(ctx, repo.q).ListSellers(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcSellers, mapSeller), nil
}
