package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetByResponse(ctx context.Context, responseID uuid.UUID) ([]models.Answer, error) {
	ctx, span := telemetry.StartSpan(ctx, "AnswerRepo.GetByResponse")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByResponse(ctx, responseID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
