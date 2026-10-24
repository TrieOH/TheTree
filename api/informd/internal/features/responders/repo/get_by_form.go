package repo

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Responder, error) {
	rows, err := database.Queries(ctx, repo.q).GetRespondersByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapResponder), nil
}
