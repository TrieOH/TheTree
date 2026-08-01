package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesQueries.ListTemplates")
	defer span.End()
	return q.repo.ListByEdition(ctx, editionID)
}
