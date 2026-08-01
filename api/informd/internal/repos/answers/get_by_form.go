package answers

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Answer, error) {
	ctx, span := telemetry.StartSpan(ctx, "AnswerRepo.GetByFormID")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetAnswersByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
