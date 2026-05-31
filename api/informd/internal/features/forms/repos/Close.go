package repos

import (
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Close(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.Close")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).CloseForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
