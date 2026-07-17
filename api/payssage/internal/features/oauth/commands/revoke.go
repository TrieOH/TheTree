package commands

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (c *Commands) Revoke(ctx context.Context, input models.RevokeInput) error {
	ctx, span := c.tracer.Start(ctx, "Revoke")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	switch input.Flow {
	case models.OAuthFlowCollector:
		var collector *models.Collector
		collector, err = c.collectors.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if collector.OrganizationID != nil {
			var org *models.Organization
			org, err = c.orgs.GetByID(ctx, *collector.OrganizationID)
			if err != nil {
				return err
			}
			err = c.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if collector.OwnerID != ident.Sub.ID {
			return fun.ErrForbidden("insufficient permissions")
		}
		return c.collectors.Revoke(ctx, input.ID)

	case models.OAuthFlowSeller:
		var seller *models.Seller
		seller, err = c.sellers.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		var wallet *models.Wallet
		wallet, err = c.wallets.GetByID(ctx, seller.WalletID)
		if err != nil {
			return err
		}
		if wallet.OrganizationID != nil {
			var org *models.Organization
			org, err = c.orgs.GetByID(ctx, *wallet.OrganizationID)
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
		return c.sellers.Revoke(ctx, input.ID)

	default:
		return fun.ErrBadRequest("invalid flow")
	}
}
