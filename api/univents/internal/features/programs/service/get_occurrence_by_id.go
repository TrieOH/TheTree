package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	return o.Occurrences.GetOccurrenceByID(ctx, id)
}
