package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.ListTemplatesByEdition")
	defer span.End()
	rows, err := database.Queries(ctx, repo.q).ListBadgeTemplatesByEdition(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(rows, mapBadgeTemplate), nil
}
