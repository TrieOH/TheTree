package service

import (
	"context"
	"univents/errs"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) Discontinue(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusActive {
		return errs.ErrEventNotActive
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Discontinue(ctx, eventID)
}
