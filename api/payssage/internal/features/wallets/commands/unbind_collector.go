package commands

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "UnbindCollector")
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

	return c.wallets.UnbindCollector(ctx, walletID)
}
