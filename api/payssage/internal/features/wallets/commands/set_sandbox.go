package commands

import (
	"context"
	"lib/telemetry"
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

	var org *models.Organization
	if payload.OrganizationID != nil {
		org, err = c.orgs.GetByID(ctx, *payload.OrganizationID)
		if err != nil {
			return err
		}

		err = c.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return err
		}
	}

	wallet, err := c.wallets.GetByID(ctx, payload.WalletID)
	if err != nil {
		return err
	}

	err = c.checkWalletAccess(wallet, ident.Sub.ID, org)
	if err != nil {
		return err
	}

	return c.wallets.SetSandboxState(ctx, wallet.ID, payload.Sandbox)
}
