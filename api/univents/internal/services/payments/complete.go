package payments

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Complete persists the seller + provider public key on the event, after the
// OAuth callback delivered them (seller_id = payssage's credential_id). The
// seller must belong to the event's wallet — verified against Payssage, so a
// forged or foreign seller id is rejected.
func (o *Operations) Complete(ctx context.Context, eventID, sellerID uuid.UUID, publicKey string) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "PaymentsService.Complete")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	if event.PayssageWalletID == nil {
		return nil, fun.ErrBadRequest("event has no payssage wallet; connect first")
	}

	sellers, err := o.payssage.ListWalletSellers(ctx, *event.PayssageWalletID)
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
		return nil, fun.ErrBadRequest("seller does not belong to the event's wallet")
	}

	return o.events.SetPaymentsConfig(ctx, event.ID, &sellerID, event.PayssageWalletID, &publicKey)
}
