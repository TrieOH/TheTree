package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListProgramsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListProgramsByEdition")
	defer span.End()
	return q.programs.ListByEdition(ctx, editionID)
}
