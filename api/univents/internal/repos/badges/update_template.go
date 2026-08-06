package badges

import (
	"context"
	"encoding/json"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) Update(ctx context.Context, id uuid.UUID, name string, designData json.RawMessage) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.UpdateTemplate")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).UpdateBadgeTemplate(ctx, sqlc.UpdateBadgeTemplateParams{
		ID:         id,
		Name:       name,
		DesignData: designData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapBadgeTemplate(row)), nil
}
