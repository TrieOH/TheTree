package service

import (
	"context"
	"lib/telemetry"
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

	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
