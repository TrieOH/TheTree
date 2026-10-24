package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListProgramsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.ListProgramsByEdition")
	defer span.End()
	return o.Programs.ListByEdition(ctx, editionID)
}
