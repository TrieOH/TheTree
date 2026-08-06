package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListViewsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEmissionView, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeEmissionsRepo.ListViewsByEdition")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListBadgeEmissionViewsByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapBadgeEmissionViewFromEdition), nil
}
