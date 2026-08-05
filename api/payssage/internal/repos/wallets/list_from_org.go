package wallets

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListFromOrg")
	defer span.End()
	sqlcWallets, err := database.Queries(ctx, repo.q).ListOrgWallets(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcWallets, mapWallet), nil
}
