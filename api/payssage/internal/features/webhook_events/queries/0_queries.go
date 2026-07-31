package queries

import (
	"context"
	"fmt"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

type Queries struct {
	events  ports.WebhookEventRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	events ports.WebhookEventRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		events:  events,
		wallets: wallets,
		orgs:    orgs,
	}
}

func (q *Queries) checkWalletAccess(ctx context.Context, walletID, subID uuid.UUID) error {
	wallet, err := q.wallets.GetByID(ctx, walletID)
	if err != nil {
		return err
	}
	if wallet.OwnerID == subID {
		return nil
	}
	if wallet.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *wallet.OrganizationID)
		if err != nil {
			return err
		}
		if org.OwnerID == subID {
			return nil
		}
		_, err = q.orgs.GetMember(ctx, subID, org.ID)
		if err != nil {
			return fmt.Errorf("insufficient permissions: %w", err)
		}
		return nil
	}
	return fun.ErrForbidden("insufficient permissions")
}
