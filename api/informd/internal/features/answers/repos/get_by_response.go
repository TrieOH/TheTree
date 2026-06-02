package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) GetByResponse(ctx context.Context, responseID uuid.UUID) ([]models.Answer, error) {
	ctx, span := repo.tracer.Start(ctx, "AnswerRepo.GetByResponse")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByResponse(ctx, responseID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
