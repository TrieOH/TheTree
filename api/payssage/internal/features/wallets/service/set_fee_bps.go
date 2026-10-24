package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) SetFeeBPS(ctx context.Context, payload models.SetFeeBPSInput) error {
	ctx, span := telemetry.StartSpan(ctx, "SetFeeBPS")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.Wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	err = o.Authz.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.Wallets.SetFeeBPS(ctx, wallet.ID, payload.FeeBps)
}
