package repos

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (r *Repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.GetTemplateByID")
	defer span.End()
	row, err := r.q.GetBadgeTemplateByID(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}

	result := mapBadgeTemplate(row)
	return &result, nil
}
