package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Seller, error) {
	ctx, span := q.tracer.Start(ctx, "ListByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet.OwnerID == ident.Sub.ID {
		return q.sellers.ListByWallet(ctx, wallet.ID)
	}

	if wallet.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return nil, err
		}
		if org.OwnerID == ident.Sub.ID {
			return q.sellers.ListByWallet(ctx, wallet.ID)
		}
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
		if err != nil {
			return nil, err
		}
		return q.sellers.ListByWallet(ctx, wallet.ID)
	}

	return nil, nil
}
