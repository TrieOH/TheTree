package repos

import (
	"context"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (r *Repo) CreateTemplate(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.CreateTemplate")
	defer span.End()
	row, err := r.q.CreateBadgeTemplate(ctx, sqlc.CreateBadgeTemplateParams{
		EditionID:    template.EditionID,
		TicketTypeID: template.TicketTypeID,
		Name:         template.Name,
		DesignData:   template.DesignData,
	})
	if err != nil {
		return nil, r.dbe(err)
	}

	result := mapBadgeTemplate(row)
	return &result, nil
}
