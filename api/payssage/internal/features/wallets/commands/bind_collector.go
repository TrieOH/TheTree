package commands

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BindCollector")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := c.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OrganizationID != nil {
		org, err := c.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		err = c.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return err
		}
	} else if wallet.OwnerID != ident.Sub.ID {
		return fun.ErrForbidden("insufficient permissions")
	}

	return c.wallets.BindCollector(ctx, walletID, collectorID)
}
