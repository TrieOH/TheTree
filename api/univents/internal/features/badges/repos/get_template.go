package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.GetTemplateByID")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetBadgeTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapBadgeTemplate(row)), nil
}
