package badges

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.CreateTemplate")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateBadgeTemplate(ctx, sqlc.CreateBadgeTemplateParams{
		EditionID:    template.EditionID,
		TicketTypeID: template.TicketTypeID,
		Name:         template.Name,
		DesignData:   template.DesignData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapBadgeTemplate(row)), nil
}
