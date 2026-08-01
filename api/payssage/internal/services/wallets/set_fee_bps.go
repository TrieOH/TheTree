package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
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

	wallet, err := o.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.wallets.SetFeeBPS(ctx, wallet.ID, payload.FeeBps)
}
