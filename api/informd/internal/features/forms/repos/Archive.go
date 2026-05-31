package repos

import (
	"Informd/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Archive(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "FormRepo.Archive")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).ArchiveForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
