package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	return o.Occurrences.ListOccurrencesByProgram(ctx, programID)
}

func (o *Operations) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	return o.Occurrences.ListOccurrencesByEdition(ctx, editionID)
}
