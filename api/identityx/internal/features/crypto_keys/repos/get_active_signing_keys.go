package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetActiveSigningKeys")
	defer span.End()

	sqlcKeys, err := database.Queries(ctx, repo.q).GetActiveSigningKeys(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcKeys, mapToActiveSigningKey), nil
}
