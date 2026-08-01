package answers

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetByField(ctx context.Context, fieldID uuid.UUID) ([]models.Answer, error) {
	ctx, span := telemetry.StartSpan(ctx, "AnswerRepo.GetByField")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByField(ctx, &fieldID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
