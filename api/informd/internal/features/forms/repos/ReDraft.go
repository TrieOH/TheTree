package repos

import (
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.ReDraft")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).DraftForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
