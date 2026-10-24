package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) PatchOccurrence(ctx context.Context, payload models.PatchProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.PatchOccurrence")
	defer span.End()

	existing, err := o.Occurrences.GetOccurrenceByID(ctx, payload.OccurrenceID)
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

	occurrence := &models.ProgramOccurrence{
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return o.Occurrences.PatchOccurrence(ctx, payload.OccurrenceID, occurrence)
}
