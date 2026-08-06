package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListViewsByUser(ctx context.Context, userID uuid.UUID) ([]models.BadgeEmissionView, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeEmissionsRepo.ListViewsByUser")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListBadgeEmissionViewsByUser(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapBadgeEmissionView), nil
}
