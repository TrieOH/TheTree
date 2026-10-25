package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// Disconnect unlinks the seller from the event (unlink only): the seller and
// the provider public key are cleared. The platform wallet itself is
// env-configured (D6) and untouched — reconnecting just re-OAuths a seller
// onto it.
func (o *Operations) Disconnect(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Events.ClearPaymentsConfig(ctx, event.ID)
}
