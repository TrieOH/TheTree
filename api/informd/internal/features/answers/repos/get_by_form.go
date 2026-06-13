package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Answer, error) {
	ctx, span := repo.tracer.Start(ctx, "AnswerRepo.GetByFormID")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
