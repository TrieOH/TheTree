package programs

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetOccurrenceByID(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.GetOccurrenceByID")
	defer span.End()
	return o.occurrences.GetOccurrenceByID(ctx, id)
}
