package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "UnbindCollector")
	defer span.End()

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
		err = o.Authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
		if err != nil {
			return err
		}
	} else if wallet.OwnerID != ident.Sub.ID {
		return fun.ErrForbidden("insufficient permissions")
	}

	return o.Wallets.UnbindCollector(ctx, walletID)
}
