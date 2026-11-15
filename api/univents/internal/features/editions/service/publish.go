package service

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"univents/errs"
	"univents/models"
	"uuid"
)

func (o *Operations) Publish(ctx context.Context, editionID uuid.UUID) error {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if !edition.IsDraft {
		return errs.ErrEditionAlreadyPublished
	}

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	err = o.Editions.Publish(ctx, editionID)
	if err != nil {
		return err
	}

	err = o.Badges.AwardStaffBadgesForEdition(ctx, editionID)
	if err != nil {
		telemetry.Log().Error("failed to award staff badges for edition",
			slog.String("edition_id", editionID.String()),
			telemetry.Err(err))
	}

	return nil
}
