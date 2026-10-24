package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.Delete")
	defer span.End()

	existing, err := o.Signatures.GetByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Signatures.Delete(ctx, id)
}
