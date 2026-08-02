package oauth

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (o *Operations) Revoke(ctx context.Context, input models.RevokeInput) error {
	ctx, span := telemetry.StartSpan(ctx, "Revoke")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	switch input.Flow {
	case models.OAuthFlowCollector:
		var collector *models.Collector
		collector, err = o.collectors.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if collector.Provider != input.Provider {
			return fun.ErrBadRequest("provider mismatch")
		}
		if collector.OrganizationID != nil {
			var org *models.Organization
			org, err = o.orgs.GetByID(ctx, *collector.OrganizationID)
			if err != nil {
				return err
			}
			err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if collector.OwnerID != ident.Sub.ID {
			return fun.ErrForbidden("insufficient permissions")
		}
		return o.collectors.Revoke(ctx, input.ID)

	case models.OAuthFlowSeller:
		var seller *models.Seller
		seller, err = o.sellers.GetByID(ctx, input.ID)
		if err != nil {
			return err
		}
		if seller.Provider != input.Provider {
			return fun.ErrBadRequest("provider mismatch")
		}
		var wallet *models.Wallet
		wallet, err = o.wallets.GetByID(ctx, seller.WalletID)
		if err != nil {
			return err
		}
		if wallet.OrganizationID != nil {
			var org *models.Organization
			org, err = o.orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return err
			}
			err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return err
			}
		} else if wallet.OwnerID != ident.Sub.ID {
			return fun.ErrForbidden("insufficient permissions")
		}
		return o.sellers.Revoke(ctx, input.ID)

	default:
		return fun.ErrBadRequest("invalid flow")
	}
}
