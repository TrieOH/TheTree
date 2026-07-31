package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetTemplate(ctx context.Context, editionID, templateID uuid.UUID) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesQueries.GetTemplate")
	defer span.End()
	return q.repo.GetTemplateByID(ctx, templateID)
}
