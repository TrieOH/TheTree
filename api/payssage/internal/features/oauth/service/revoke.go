package service

import (
	"context"
	"payssage/errs"
	"payssage/models"
	idx "sdk/identityx"
)

func (o *Operations) Revoke(ctx context.Context, input models.RevokeInput) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	switch input.Flow {
	case models.OAuthFlowCollector:
		var collector *models.Collector
		collector, err = o.Collectors.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if collector.Provider != input.Provider {
			return errs.ErrProviderMismatch
		}
		if collector.OrganizationID != nil {
			var org *models.Organization
			org, err = o.Orgs.GetByID(ctx, *collector.OrganizationID)
			if err != nil {
				return err
			}
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if collector.OwnerID != ident.Sub.ID {
			return errs.ErrNotResourceOwner
		}
		return o.Collectors.Revoke(ctx, input.ID)

	case models.OAuthFlowSeller:
		var seller *models.Seller
		seller, err = o.Sellers.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if seller.Provider != input.Provider {
			return errs.ErrProviderMismatch
		}
		var wallet *models.Wallet
		wallet, err = o.Wallets.GetByID(ctx, seller.WalletID)
		if err != nil {
			return err
		}
		if wallet.OrganizationID != nil {
			var org *models.Organization
			org, err = o.Orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return err
			}
			_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if wallet.OwnerID != ident.Sub.ID {
			return errs.ErrNotResourceOwner
		}
		return o.Sellers.Revoke(ctx, input.ID)

	default:
		return errs.ErrInvalidFlow
	}
}
