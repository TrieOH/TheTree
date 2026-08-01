package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetTemplate(ctx context.Context, templateID uuid.UUID) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesQueries.GetTemplate")
	defer span.End()
	return o.repo.GetByID(ctx, templateID)
}
