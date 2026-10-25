package service

import (
	"context"
	"payssage/errs"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.Wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OrganizationID != nil {
		org, err := o.Orgs.GetByID(ctx, *wallet.OrganizationID)
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

	return o.Wallets.BindCollector(ctx, walletID, collectorID)
}
