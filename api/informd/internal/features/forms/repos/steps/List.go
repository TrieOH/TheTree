package steps

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *stepRepo) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	ctx, span := repo.tracer.Start(ctx, "FormRepo.List")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).ListStepsByFormID(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForm, mapStep), nil
}
