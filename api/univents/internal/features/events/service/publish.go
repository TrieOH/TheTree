package service

import (
	"context"
	"univents/errs"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) Publish(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusDraft {
		return errs.ErrEventNotDraft
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Publish(ctx, eventID)
}
