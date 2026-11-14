package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	return o.Occurrences.GetOccurrenceByID(ctx, id)
}
