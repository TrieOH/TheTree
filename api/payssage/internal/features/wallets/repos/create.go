package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Wallet) (*models.Wallet, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcWallet, err := database.Queries(ctx, repo.q).CreateWallet(ctx, sqlc.CreateWalletParams{
		OwnerID:        toCreate.OwnerID,
		OrganizationID: toCreate.OrganizationID,
		Name:           toCreate.Name,
		Sandbox:        toCreate.Sandbox,
		FeeBps:         toCreate.FeeBps,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapWallet(sqlcWallet)), nil
}
