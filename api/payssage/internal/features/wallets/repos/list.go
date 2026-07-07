package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := repo.tracer.Start(ctx, "List")
	defer span.End()
	sqlcWallets, err := database.Queries(ctx, repo.q).ListWallets(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcWallets, mapWallet), nil
}
