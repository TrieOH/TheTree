package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error) {
	sqlcSellers, err := database.Queries(ctx, repo.q).ListSellersByWallet(ctx, walletID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcSellers, mapSeller), nil
}
