package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return nil, err
		}

		err = q.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleMember)
		if err != nil {
			return nil, err
		}

		return wallet, nil
	}

	if !wallet.OwnedBy(ident.Sub.ID) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return wallet, nil
}
