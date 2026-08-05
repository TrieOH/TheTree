package wallets

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BindCollector")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	wallet, err := o.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet.OrganizationID != nil {
		org, err := o.orgs.GetByID(ctx, *wallet.OrganizationID)
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

	return o.wallets.BindCollector(ctx, walletID, collectorID)
}
