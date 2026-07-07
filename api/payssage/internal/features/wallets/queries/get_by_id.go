package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	ctx, span := q.tracer.Start(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet.OwnerID != ident.Sub.ID && wallet.OrganizationID == nil {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	if wallet.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return nil, err
		}
		if ident.Sub.ID != org.OwnerID {
			_, err = q.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return wallet, nil
}
