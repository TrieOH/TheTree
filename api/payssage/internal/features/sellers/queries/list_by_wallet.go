package queries

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return q.sellers.ListByWallet(ctx, wallet.ID)
}
