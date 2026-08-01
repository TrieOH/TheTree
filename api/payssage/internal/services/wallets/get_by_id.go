package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := o.wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
