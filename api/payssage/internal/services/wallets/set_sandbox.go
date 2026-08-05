package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) SetSandbox(ctx context.Context, payload models.SetSandboxInput) error {
	ctx, span := telemetry.StartSpan(ctx, "SetSandbox")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return o.wallets.SetSandboxState(ctx, wallet.ID, payload.Sandbox)
}
