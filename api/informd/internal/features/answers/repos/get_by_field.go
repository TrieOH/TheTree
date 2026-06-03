package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) GetByField(ctx context.Context, fieldID uuid.UUID) ([]models.Answer, error) {
	ctx, span := repo.tracer.Start(ctx, "AnswerRepo.GetByField")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByField(ctx, &fieldID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
