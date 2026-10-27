package service

import (
	"context"
	"univents/models"
)

func (o *Operations) CreateOccurrence(ctx context.Context, payload models.CreateProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	program, err := o.Programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, program.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
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

	return o.Occurrences.CreateOccurrence(ctx, occurrence)
}
