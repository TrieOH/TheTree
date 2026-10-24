package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.DeleteOccurrence")
	defer span.End()

	existing, err := o.Occurrences.GetOccurrenceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Occurrences.DeleteOccurrence(ctx, id)
}
