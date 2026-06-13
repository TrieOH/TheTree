package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) GetByFormID(ctx context.Context, formID uuid.UUID) ([]models.Responder, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponderRepo.GetByFormID")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).GetRespondersByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapResponder), nil
}
