package service

import (
	"context"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Publish(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusDraft {
		return fun.ErrBadRequest("cannot publish non draft event")
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Publish(ctx, eventID)
}
