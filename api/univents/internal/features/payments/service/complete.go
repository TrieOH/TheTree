package service

import (
	"context"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Complete persists the seller + provider public key on the event, after the
// OAuth callback delivered them (seller_id = payssage's credential_id). The
// seller must belong to the platform wallet (env-configured, D6) — verified
// against Payssage, so a forged or foreign seller id is rejected.
func (o *Operations) Complete(ctx context.Context, eventID, sellerID uuid.UUID, publicKey string) (*models.Event, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	sellers, err := o.Payssage.ListWalletSellers(ctx, o.WalletID)
	if err != nil {
		return nil, err
	}
	belongs := false
	for _, s := range sellers {
		if s.ID == sellerID {
			belongs = true
			break
		}
	}
	if !belongs {
		return nil, fun.ErrBadRequest("seller does not belong to the platform payssage wallet")
	}

	return o.Events.SetPaymentsConfig(ctx, event.ID, &sellerID, &publicKey)
}
