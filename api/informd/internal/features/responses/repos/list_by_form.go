package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListByForm(ctx context.Context, formID uuid.UUID) ([]models.Response, error) {
	ctx, span := telemetry.StartSpan(ctx, "ResponseRepo.ListByForm")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListResponsesByForm(ctx, formID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapResponse), nil
}
