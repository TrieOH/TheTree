package service

import (
	"context"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Discontinue(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusActive {
		return fun.ErrBadRequest("cannot discontinue non active event")
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Discontinue(ctx, eventID)
}
