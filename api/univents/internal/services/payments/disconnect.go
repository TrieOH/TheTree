package payments

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// Disconnect unlinks the seller from the event (unlink only): the seller and
// the provider public key are cleared, the wallet is kept — it is the event's
// permanent payment container, and reconnecting reuses it.
func (o *Operations) Disconnect(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "PaymentsService.Disconnect")
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

	return o.events.ClearPaymentsConfig(ctx, event.ID)
}
