package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Answer, error) {
	rows, err := database.Queries(ctx, repo.q).GetAnswersByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapAnswer), nil
}
