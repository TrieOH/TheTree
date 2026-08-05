package programs

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) PatchOccurrence(ctx context.Context, payload models.PatchProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.PatchOccurrence")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := o.occurrences.GetOccurrenceByID(ctx, payload.OccurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	occurrence := &models.ProgramOccurrence{
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return o.occurrences.PatchOccurrence(ctx, payload.OccurrenceID, occurrence)
}
