package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

func (o *Operations) Publish(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Publish")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if !edition.IsDraft {
		return fun.ErrBadRequest("edition is already published")
	}

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
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
			zap.String("edition_id", editionID.String()),
			zap.Error(err))
	}

	return nil
}
