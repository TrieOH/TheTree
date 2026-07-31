package commands

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"
)

func (c *Commands) SetSandbox(ctx context.Context, payload models.SetSandboxInput) error {
	ctx, span := telemetry.StartSpan(ctx, "SetSandbox")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := c.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleAdmin)
	if err != nil {
		return err
	}

	return c.wallets.SetSandboxState(ctx, wallet.ID, payload.Sandbox)
}
