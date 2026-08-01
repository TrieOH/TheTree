package programs

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProgramByID(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.GetProgramByID")
	defer span.End()
	return o.programs.GetByID(ctx, id)
}
