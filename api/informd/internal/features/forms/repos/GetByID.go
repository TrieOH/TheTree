package repos

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.GetByID")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).GetFormByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
