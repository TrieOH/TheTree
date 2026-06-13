package repos

import (
	"context"

	"Informd/models"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.Open")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).OpenForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
