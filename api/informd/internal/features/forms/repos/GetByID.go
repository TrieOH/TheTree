package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormRepo.GetByID")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).GetFormByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapForm(sqlcForm)), nil
}
