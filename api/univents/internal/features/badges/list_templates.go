package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesQueries.ListTemplates")
	defer span.End()
	return o.repo.ListByEdition(ctx, editionID)
}
