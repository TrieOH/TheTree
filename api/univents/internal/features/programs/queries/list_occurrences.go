package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListOccurrencesByProgram(ctx context.Context, programID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListOccurrencesByProgram")
	defer span.End()
	return q.occurrences.ListOccurrencesByProgram(ctx, programID)
}

func (q *Queries) ListOccurrencesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListOccurrencesByEdition")
	defer span.End()
	return q.occurrences.ListOccurrencesByEdition(ctx, editionID)
}
