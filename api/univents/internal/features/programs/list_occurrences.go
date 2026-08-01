package programs

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListOccurrencesByProgram")
	defer span.End()
	return o.occurrences.ListOccurrencesByProgram(ctx, programID)
}

func (o *Operations) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListOccurrencesByEdition")
	defer span.End()
	return o.occurrences.ListOccurrencesByEdition(ctx, editionID)
}
