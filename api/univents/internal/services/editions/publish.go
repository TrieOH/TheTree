package editions

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Publish(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.Publish")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if !edition.IsDraft {
		return fun.ErrBadRequest("edition is already published")
	}

	event, err := o.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.editions.Publish(ctx, editionID)
}
