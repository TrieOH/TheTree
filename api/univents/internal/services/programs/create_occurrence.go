package programs

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) CreateOccurrence(ctx context.Context, payload models.CreateProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.CreateOccurrence")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	program, err := o.programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, program.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	occurrence := &models.ProgramOccurrence{
		ProgramID:   payload.ProgramID,
		EditionID:   program.EditionID,
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return o.occurrences.CreateOccurrence(ctx, occurrence)
}
