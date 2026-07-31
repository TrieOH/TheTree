package repos

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *Repo) ListMineArchived(ctx context.Context, userID uuid.UUID) ([]models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMineArchived")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).ListMyArchivedForms(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForms, mapForm), nil
}
